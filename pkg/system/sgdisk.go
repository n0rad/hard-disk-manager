package system

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
)

type SgDisk struct {
	exec runner.Exec
}

func (s SgDisk) Clear(device BlockDevice) error {
	std, err := s.exec.ExecGetStd("sgdisk", "-og", device.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("std", std).WithField("disk", device), "Failed to clear partition table")
	}
	return nil
}

func (s SgDisk) CreatePartition(device BlockDevice, label string) error {
	std, err := s.exec.ExecGetStd("sgdisk", "-n", "1:0:0", "-t", "1:CA7D7CCB-63ED-4C53-861C-1742536059CC", "-c", "1:"+label, device.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("std", std).WithField("disk", device), "Failed to create single partition")
	}
	return nil
}