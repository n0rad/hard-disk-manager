package system

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
)

func (b BlockDevice) ClearPartitionTable() error {
	std, err := b.exec.ExecGetStd("sgdisk", "-og", b.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("std", std), "Failed to clear partition table")
	}
	return nil
}

func (b BlockDevice) CreateSinglePartition(label string) error {
	std, err := b.exec.ExecGetStd("sgdisk", "-n", "1:0:0", "-t", "1:"+luksPartitionCode, "-c", "1:"+label, b.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("std", std), "Failed to create single partition")
	}
	return nil
}