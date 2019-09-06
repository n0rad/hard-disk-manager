package agent

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
)

type DiskManager struct {
	PassService password.Service
	Path        string

	handlers []Handler
	stop     chan struct{}
	//serialJobs chan
}

func (d *DiskManager) Init() {
	d.handlers = append(d.handlers, &HandlerDb{})
	d.handlers = append(d.handlers, &HandlerAdd{})

	for _, v := range d.handlers {
		v.Init(d)
	}
}

func (d *DiskManager) Start() error {
	logs.WithField("path", d.Path).Info("New disk manager")
	d.stop = make(chan struct{})

	for _, v := range d.handlers {
		v.Start()
	}

	//<-d.stop

	//for _, v := range d.handlers {
	//	v.Stop()
	//}
	return nil
}

func (d *DiskManager) Stop(e error) {
	close(d.stop)
}

//func (d *DiskManager) AddChildrenEvent(event hdm.BlockDeviceEvent) {
//	logs.WithField("event", event).Info("Children event")
//}

// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
