package handlers

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

// hold handlers for a block system (disk, fs or path)
type BlockManager struct {
	PassService    *password.Service
	ManagerService *ManagersService
	Path           string
	FStype         string
	Type           string

	server system.Server
	config hdm.Config // TODO that sux and blockManager should be specialized

	handlers   []Handler
	stop       chan struct{}
	serialJobs chan func() // reduce pressure on disk
}

func (d *BlockManager) Init() error {
	logs.WithField("path", d.Path).Info("New block device manager")
	if err := d.server.Init(); err != nil {
		return errs.WithE(err, "Failed to init empty server")
	}

	for _, handler := range handlers {
		if handler.filter.Match(HandlerFilter{Type: d.Type, FSType: d.FStype}) {
			handler := handler.new()
			logs.WithField("handler", handler.HandlerName()).WithField("path", d.Path).Debug("Register handler")
			d.handlers = append(d.handlers, handler)

			// TODO load configuration for handler
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
		logs.WithField("path", d.Path).WithField("handler", v.HandlerName()).Trace("Starting handler")
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
}

func (d *BlockManager) Notify(level logs.Level, fields data.Fields, message string) {
	logs.WithFields(fields).Info(message)
}

func (d *BlockManager) Event() {

	// add
	// remove
	// changed
	// mount
}

//func (d *BlockManager) AddChildrenEvent(event hdm.BlockDeviceEvent) {
//	logs.WithField("event", event).Info("Children event")
//}

// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
