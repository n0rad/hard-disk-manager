package app

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/agent"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
	"sync"
)

type Agent struct {
	PassService password.Service

	server       system.Server
	udevConn     *netlink.UEventConn
	stop         chan struct{}
	diskManagers map[string]agent.DiskManager
	disksMutex   sync.Mutex
}

func (a *Agent) Start() error {
	a.diskManagers = make(map[string]agent.DiskManager)
	a.stop = make(chan struct{})

	a.udevConn = new(netlink.UEventConn)
	defer a.udevConn.Close()
	if err := a.udevConn.Connect(netlink.UdevEvent); err != nil {
		return errs.WithE(err, "Unable to connect to Netlink Kobject UEvent socket")
	}

	if err := a.server.Init(); err != nil {
		return errs.WithE(err, "Failed to init empty server")
	}

	if err := a.addCurrentBlockDevices(); err != nil {
		a.Stop(err)
		return errs.WithE(err, "Cannot add current block devices after watching events")
	}

	// TODO you can lose events between addCurrent and watch but watch is blocking

	a.watchUdevBlockEvents()

	logs.Info("Stop Agent")

	// cleanup
	a.disksMutex.Lock()
	defer a.disksMutex.Unlock()
	//for _, v := range a.diskManagers {
	//v.Stop(nil)
	//}

	return nil
}

func (a *Agent) Stop(e error) {
	logs.Info("Stop Agent1")
	close(a.stop)
	logs.Info("Stop Agent2")
}

///////////////////////

type BlockDeviceEvent struct {
	Action netlink.KObjAction
	Path   string
	Type   string
}

func (a *Agent) addDisk(path string) {
	if _, ok := a.diskManagers[path]; !ok {
		dm := agent.DiskManager{
			Path:        path,
			PassService: a.PassService,
		}
		dm.Init()
		a.diskManagers[path] = dm
		go dm.Start()
		//if err := start; err != nil {
		//	logs.WithE(err).Error("Failed to start agent service")
		//} else {
		//}
	} else {
		logs.WithField("path", path).Warn("Cannot add disk, already exists")
	}
}

func (a *Agent) removeDisk(path string) {
	if diskManager, ok := a.diskManagers[path]; ok {
		diskManager.Stop(nil)
		delete(a.diskManagers, path)
	} else {
		logs.WithField("path", path).Warn("Cannot remove disk, not found")
	}
}

func (a *Agent) addCurrentBlockDevices() error {
	blockDevices, err := a.server.ListFlatBlockDevices()
	if err != nil {
		return errs.WithE(err, "Failed to list current block devices")
	}
	for _, v := range blockDevices {
		a.handleEvent(BlockDeviceEvent{
			Action: netlink.ADD,
			Type:   v.Type,
			Path:   v.Path,
		})
	}
	return nil
}

func (a *Agent) handleEvent(event BlockDeviceEvent) {
	a.disksMutex.Lock()
	defer a.disksMutex.Unlock()

	switch event.Type {
	case "disk":
		switch event.Action {
		case "add":
			a.addDisk(event.Path)
		case "remove":
			a.removeDisk(event.Path)
		case "change":
			a.removeDisk(event.Path)
			a.addDisk(event.Path)
		default:
			logs.WithField("event", event).Warn("Unknown udev event type")
		}
	case "part", "partition":
		logs.WithField("event", event).Info("Children event")
		//if manager, ok := a.diskManagers[event.Path]; ok {
		//	manager.AddChildrenEvent(event)
		//} else {
		//	logs.WithField("event", event).Warn("Disk not found to add event")
		//	// disk not found to add partition
		//}
	default:
		logs.WithField("event", event).Warn("Unknown event type")
	}
}

func (a *Agent) watchUdevBlockEvents() {
	matcher := netlink.RuleDefinitions{
		Rules: []netlink.RuleDefinition{
			{
				//Action: "",
				Env: map[string]string{
					"SUBSYSTEM": "block",
					//"DEVTYPE":   "disk",
					//  ID_BUS=ata
					//	ID_TYPE=disk
				},
			},
		},
	}

	queue := make(chan netlink.UEvent)
	defer close(queue)
	errors := make(chan error)
	defer close(errors)
	quitMonitor := a.udevConn.Monitor(queue, errors, &matcher)
	for {
		select {
		case uevent := <-queue:
			logs.WithField("uevent", uevent).Trace("Received udev event")
			a.handleEvent(BlockDeviceEvent{
				Action: uevent.Action,
				Path:   uevent.Env["DEVNAME"],
				Type:   uevent.Env["DEVTYPE"],
			})
		case err := <-errors:
			logs.WithE(err).Warn("Received error for udev watcher")
		case <-a.stop:
			close(quitMonitor)
			return
		}
	}
}
