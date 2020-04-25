package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type DiskManager struct {
	BlockManager

	serialJobs chan job
}

func (m *DiskManager) Init(lsblk system.Lsblk, disk string) error {
	block, err := lsblk.GetBlockDevice(disk)
	if err != nil {
		return errs.WithE(err, "Failed to get block device to init manager")
	}

	if block.Type != "disk" {
		return errs.WithF(data.WithField("disk", disk), "Not a disk device")
	}

	m.BlockManager.Init(nil, block)
	m.serialJobs = make(chan job, 5)
	return nil
}

func (m *DiskManager) Start() error {
	if err := m.preStart(); err != nil {
		return err
	}

	m.handleSerialJobs()

	return m.postStart()
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

func (m *DiskManager) handleSerialJobs() {
	for {
		select {
		case job := <-m.serialJobs:
			res := job.f()
			job.done <- res
			// TODO close channel ?
			close(job.done)
		case <-m.stop:
			return
		}
	}
}
