package agent

import (
	"github.com/n0rad/go-erlog/logs"
)

type DiskManager struct {
	path string
	handlers []Handler

	//serialJobs chan
}

func (d *DiskManager) Stop() {
	logs.WithField("path", d.path).Info("Stop disk manager")
}

//func (d *DiskManager) AddChildrenEvent(event hdm.BlockDeviceEvent) {
//	logs.WithField("event", event).Info("Children event")
//}

func NewDiskManager(path string) DiskManager {
	handlers := []Handler {
		&HandlerDb{},
	}

	for _, v := range handlers {
		v.Init(path)
	}


	for _, v := range handlers {
		v.Start()
	}





	logs.WithField("path", path).Info("New disk manager")

	//disk, err := server.ScanDisk(path)
	//if err != nil {
	//	return DiskManager{}, errs.WithEF(err, data.WithField("path", path), "Failed to scan disk")
	//}
	//
	//// store disk in db
	////
	manager := DiskManager{
		path: path,
		//disk: disk,
	}

	return manager
}


// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
