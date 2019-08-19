package agent

import "github.com/n0rad/go-erlog/logs"

type DiskManager struct {
	path string
	// events
}

func (d *DiskManager) Stop() {
	logs.WithField("path", d.path).Info("Stop disk manager")
}

func NewDiskManager(path string) DiskManager {
	logs.WithField("path", path).Info("New disk manager")

	// get info
	// store disk in db
	//
	return DiskManager{path: path}
}

// register in DB
// timer check health
// date test health
// inotify on files
// timer/initify index files
// prepare if disk is empty
//

// watch disk events -> add/remove disks
// watch files events -> run sync
// scan for backup -> timing sync
// timer rsync
// date health
