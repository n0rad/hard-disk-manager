package trash

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

func Prepare(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)
	//label := selector.Label
	selector.Label = ""

	return hdm.HDM.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		if disk.HasChildren() {
			return errs.WithF(fields, "Cannot prepare, disk has partitions")
		}

		//password, err := utils.AskPasswordWithConfirmation(true)
		//if err != nil {
		//	return errs.WithE(err, "Failed to get password")
		//}
		//
		//return disk.Prepare(label, string(password))
		return nil
	})
}

func Erase(selector system.DisksSelector) error {
	return hdm.HDM.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		return disk.Erase()
	})
}

func Add(selector system.DisksSelector) error {
	//password, err := utils.AskPasswordWithConfirmation(false)
	//if err != nil {
	//	return errs.WithE(err, "Failed to get password")
	//}

	return hdm.HDM.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		//return disk.AddBlockDevice()
		return nil
	})
}

func Remove(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)

	disk, err := hdm.HDM.Servers.GetDisk(selector)
	if err != nil {
		return err
	}
	if disk == nil {
		return errs.WithF(fields, "Disk not found")
	}

	return disk.Remove()
}
