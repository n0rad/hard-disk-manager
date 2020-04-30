package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
)

type DiskManager struct {
	BlockManager

	serialJobs chan job
}

func (m *DiskManager) Init(lsblk *system.Lsblk, disk string, udev *system.UdevService) error {
	block, err := lsblk.GetBlockDevice(disk)
	if err != nil {
		return errs.WithE(err, "Failed to get block device to init manager")
	}

	if block.Type != "disk" {
		return errs.WithF(data.WithField("disk", disk), "Not a disk device")
	}

	m.BlockManager.Init(nil, lsblk, disk, udev)
	m.serialJobs = make(chan job, 5)
	return nil
}


func (m *DiskManager) Start() error {
	udevChan := m.udev.Watch(m.block.Path)
	defer m.udev.Unwatch(udevChan)

	if err := m.preStart(); err != nil {
		return err
	}

	for {
		select {
		case job := <-m.serialJobs:
			res := job.f()
			job.done <- res
			close(job.done)

		case event := <-udevChan:
			if event.Action == netlink.ADD {
				m.HandleEvent(Add)
			} else if event.Action == netlink.REMOVE {
				m.HandleEvent(Remove)
			}

		case <-m.stop:
			return m.postStart()
		}
	}
}

////////////////////////

type job struct {
	f    func() interface{}
	done chan interface{}
}

func (m *DiskManager) runSerialJob(f func() interface{}) <-chan interface{} {
	job := job{
		f:    f,
		done: make(chan interface{}, 1),
	}
	m.serialJobs <- job
	return job.done
}
