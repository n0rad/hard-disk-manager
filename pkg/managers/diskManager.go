package managers

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/managers/block"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)


///////

type DiskManager struct {
	block.Manager

	serialJobs chan func() // reduce pressure on disk
}

func (m *DiskManager) Init(lsblk system.Lsblk, disk string) error {
	block, err := lsblk.GetBlockDevice(disk)
	if err != nil {
		return errs.WithE(err, "Failed to get block device to init manager")
	}

	if block.Type != "disk" {
		return errs.WithF(data.WithField("disk", disk), "Not a disk device")
	}

	m.Manager.Init(block, &hdm.HDM) // TODO hdm init
	return nil
}

func (m *DiskManager) Start() error {
	if err := m.Manager.Start(); err != nil {
		return err
	}
	return nil
}

func (m *DiskManager) Stop(err error) {
	m.Manager.Stop(err)
}
