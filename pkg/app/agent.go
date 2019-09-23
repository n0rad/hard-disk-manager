package app

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/handlers"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
	"sync"
)

type Agent struct {
	PassService *password.Service

	server       system.Server
	udevConn     *netlink.UEventConn
	stop         chan struct{}
	diskManagers map[string]handlers.BlockDeviceManager
	disksMutex   sync.Mutex
}

func (a *Agent) Start() error {
	a.diskManagers = make(map[string]handlers.BlockDeviceManager)
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
	close(a.stop)
}

///////////////////////

type BlockDeviceEvent struct {
	Action netlink.KObjAction
	Path   string
	Type   string
	FSType string
}

func (a *Agent) addDisk(event BlockDeviceEvent) {
	if _, ok := a.diskManagers[event.Path]; !ok {
		manager := handlers.BlockDeviceManager{
			Path:        event.Path,
			Type: event.Type,
			FStype: event.FSType,
			PassService: a.PassService,
		}

		if err := manager.Init(); err != nil {
			logs.Warn("Failed to init blockdevice manager")
		}
		a.diskManagers[event.Path] = manager
		go manager.Start()
		//if err := start; err != nil {
		//	logs.WithE(err).Error("Failed to start agent service")
		//} else {
		//}
	} else {
		logs.WithField("path", event.Path).Warn("Cannot add disk, already exists")
	}
}

func (a *Agent) removeDisk(event BlockDeviceEvent) {
	if diskManager, ok := a.diskManagers[event.Path]; ok {
		diskManager.Stop(nil)
		delete(a.diskManagers, event.Path)
	} else {
		logs.WithField("path", event.Path).Warn("Cannot remove disk, not found")
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
			FSType: v.Fstype,
		})
	}
	return nil
}

func (a *Agent) handleEvent(event BlockDeviceEvent) {
	a.disksMutex.Lock()
	defer a.disksMutex.Unlock()

	switch event.Action {
	case "add":
		a.addDisk(event)
	case "remove":
		a.removeDisk(event)
	case "change":
		a.removeDisk(event)
		a.addDisk(event)
	default:
		logs.WithField("event", event).Warn("Unknown udev event action")
	}
}

//switch event.Type {
//case "disk":
//case "part", "partition":
//	logs.WithField("event", event).Info("Children event")
//if manager, ok := a.diskManagers[event.Path]; ok {
//	manager.AddChildrenEvent(event)
//} else {
//	logs.WithField("event", event).Warn("Disk not found to add event")
//	// disk not found to add partition
//}
//default:
//	logs.WithField("event", event).Warn("Unknown event type")
//}

func (a *Agent) watchUdevBlockEvents() {
	matcher := netlink.RuleDefinitions{
		Rules: []netlink.RuleDefinition{
			{
				Env: map[string]string{
					"SUBSYSTEM": "block",
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
				FSType: uevent.Env["ID_FS_TYPE"],
			})
		case err := <-errors:
			logs.WithE(err).Warn("Received error for udev watcher")
		case <-a.stop:
			close(quitMonitor)
			return
		}
	}
}
