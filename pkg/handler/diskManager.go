package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type DiskManager struct {
	BlockManager

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

	m.BlockManager.Init(nil, block)
	return nil
}
