package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/agent"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
	"sync"
	"time"
)

type Agent struct {
	server     system.Server
	udevConn   *netlink.UEventConn
	stop       chan struct{}
	disks      map[string]agent.DiskManager
	disksMutex sync.Mutex
}

func (a *Agent) Start() error {
	a.disks = make(map[string]agent.DiskManager)
	a.stop = make(chan struct{})
	a.udevConn = new(netlink.UEventConn)
	if err := a.udevConn.Connect(netlink.UdevEvent); err != nil {
		return errs.WithE(err, "Unable to connect to Netlink Kobject UEvent socket")
	}

	if err := a.server.Init(); err != nil {
		return errs.WithE(err, "Failed to init empty server")
	}

	go a.periodicDiskWatch(func(disks []string) {
		a.disksMutex.Lock()
		for _, v := range disks {
			if _, ok := a.disks[v]; !ok {
				a.disks[v] = agent.NewDiskManager(v)
			}
		}

		for k, v := range a.disks {
			found := false
			for _, disk := range disks {
				if k == disk {
					found = true
					break
				}
			}
			if !found {
				logs.WithField("dev", k).Warn("Found disk by full scan that should not exists")
				v.Stop()
				delete(a.disks, k)
			}
		}
		a.disksMutex.Unlock()
	})

	go a.watchUdevDiskEvents(func(event netlink.UEvent) {
		a.disksMutex.Lock()
		switch event.Action {
		case "add":
			if _, ok := a.disks[event.Env["DEVNAME"]]; !ok {
				a.disks[event.Env["DEVNAME"]] = agent.NewDiskManager(event.Env["DEVNAME"])
			} else {
				logs.WithField("dev", event.Env["DEVNAME"]).Warn("Received add on already existing disk")
			}
		case "remove":
			if diskManager, ok := a.disks[event.Env["DEVNAME"]]; ok {
				diskManager.Stop()
				delete(a.disks, event.Env["DEVNAME"])
			} else {
				logs.WithField("dev", event.Env["DEVNAME"]).Warn("Received remove on not found disk")
			}
		case "change":
			if diskManager, ok := a.disks[event.Env["DEVNAME"]]; ok {
				diskManager.Stop()
				delete(a.disks, event.Env["DEVNAME"])
			}
			a.disks[event.Env["DEVNAME"]] = agent.NewDiskManager(event.Env["DEVNAME"])
		default:
			logs.WithField("event", event).Warn("Unknown udev event type")
		}
		a.disksMutex.Unlock()
	})
	return nil
}

func (a *Agent) Stop() {
	a.stop <- struct{}{}
	_ = a.udevConn.Close()
}

func (a *Agent) periodicDiskWatch(handler func(disks []string)) {
	updateDisksList := func() {
		disks, err := a.server.ListDisks()
		if err != nil {
			logs.WithE(err).Warn("Failed to scan disks")
		}
		handler(disks)
	}

	updateDisksList()

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			updateDisksList()
		case <-a.stop:
			ticker.Stop()
			return
		}
	}
}

func (a *Agent) watchUdevDiskEvents(handler func(event netlink.UEvent)) {
	matcher := netlink.RuleDefinitions{
		Rules: []netlink.RuleDefinition{
			{
				//Action: "",
				Env: map[string]string{
					"SUBSYSTEM": "block",
					"DEVTYPE":   "disk",
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
	quit := a.udevConn.Monitor(queue, errors, &matcher)
	for {
		select {
		case uevent := <-queue:
			logs.WithField("uevent", uevent).Trace("Received udev event")
			handler(uevent)
		case err := <-errors:
			logs.WithE(err).Warn("Received error for udev watcher")
		case <-a.stop:
			quit <- struct{}{}
			return
		}
	}
}
