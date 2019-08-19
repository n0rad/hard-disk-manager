package hdm

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	system "github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
)

func (hdm *Hdm) Prepare(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)
	label := selector.Label
	selector.Label = ""

	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		if disk.HasChildren() {
			return errs.WithF(fields, "Cannot prepare, disk has partitions")
		}

		password, err := utils.AskPasswordWithConfirmation(true)
		if err != nil {
			return errs.WithE(err, "Failed to get password")
		}

		return disk.Prepare(label, password)
	})
}

func (hdm *Hdm) Erase(selector system.DisksSelector) error {
	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		return disk.Erase()
	})
}

func (hdm *Hdm) Add(selector system.DisksSelector) error {
	password, err := utils.AskPasswordWithConfirmation(false)
	if err != nil {
		return errs.WithE(err, "Failed to get password")
	}

	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		return disk.Add(password)
	})
}

func (hdm *Hdm) Remove(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)

	disk, err := hdm.Servers.GetDisk(selector)
	if err != nil {
		return err
	}
	if disk == nil {
		return errs.WithF(fields, "Disk not found")
	}

	return disk.Remove()
}