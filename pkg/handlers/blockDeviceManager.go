package handlers

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

// disk
// part
// loop
// crypt
// lvm
// raid1
// md
// rom
// dmraid

// part, lvm, crypt, dmraid, mpath, path, dm, loop, md, linear, raid0, raid1, raid4, raid5, raid10, multipath, disk, tape, printer, processor, worm, rom, scanner, mo-disk, changer, comm, raid, enclosure, rbc, osd, and no-lun

type BlockDeviceManager struct {
	PassService *password.Service
	Path        string
	FStype      string
	Type        string

	server system.Server

	handlers []Handler
	stop     chan struct{}
	//serialJobs chan
}

func (d *BlockDeviceManager) Init() error {
	logs.WithField("path", d.Path).Info("New block device manager")
	if err := d.server.Init(); err != nil {
		return errs.WithE(err, "Failed to init empty server")
	}

	for _, handler := range handlers {
		if handler.filter.Match(HandlerFilter{Type: d.Type, FSType: d.FStype}) {
			handler := handler.new()
			logs.WithField("handler", handler.Name()).WithField("path", d.Path).Debug("Register handler")
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

func (d *BlockDeviceManager) Start() error {
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

func (d *BlockDeviceManager) Stop(e error) {
	close(d.stop)
}

func (d *BlockDeviceManager) Notify(level logs.Level, fields data.Fields, message string) {
	logs.WithFields(fields).Info(message)
}

func (d *BlockDeviceManager) Event() {

	// add
	// remove
	// changed
	// mount
}

//func (d *BlockDeviceManager) AddChildrenEvent(event hdm.BlockDeviceEvent) {
//	logs.WithField("event", event).Info("Children event")
//}

// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
