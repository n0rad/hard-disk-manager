package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

// hold manager for a block system (disk, fs or path)
type BlockManager struct {
	PassService    *password.Service
	ManagerService *hdm.ManagersService

	Path   string
	FStype string
	Type   string

	// TODO
	BlockDevice system.BlockDevice

	configPath string
	config     hdm.Config // TODO that sux and blockManager should be specialized

	handlers   []Handler
	stop       chan struct{}
	serialJobs chan func() // reduce pressure on disk

	childrens []BlockManager
}

func (d *BlockManager) Init() error {
	logs.WithField("path", d.Path).Info("new block device manager")

	for _, handler := range handlers {
		if handler.filter.Match(HandlerFilter{Type: d.Type, FSType: d.FStype}) {
			handler := handler.new()
			logs.WithField("manager", handler.HandlerName()).WithField("path", d.Path).Debug("register manager")
			d.handlers = append(d.handlers, handler)

			// TODO load configuration for manager
			// TODO if disabled, remove
		}
	}

	for _, v := range d.handlers {
		v.Init(d)
	}
	return nil
}

func (d *BlockManager) Start() error {
	d.stop = make(chan struct{})

	for _, v := range d.handlers {
		logs.WithField("path", d.Path).WithField("manager", v.HandlerName()).Trace("Starting manager")
		go v.Start()
	}

	d.handleSerialJobs()

	for _, v := range d.handlers {
		v.Stop()
	}
	return nil
}

func (d *BlockManager) handleSerialJobs() {
	for {
		select {
		case job := <-d.serialJobs:
			job()
		case <-d.stop:
			return
		}
	}
}

func (d *BlockManager) Stop(e error) {
	close(d.stop)
	// TODO waitgroup
}

//func (d *BlockManager) Notify(level logs.Level, fields data.Fields, message string) {
//	logs.WithFields(fields).Info(message)
//}

//func (d *BlockManager) AddChildrenEvent(event hdm.BlockDeviceEvent) {
//	logs.WithField("event", event).Info("children event")
//}

// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
