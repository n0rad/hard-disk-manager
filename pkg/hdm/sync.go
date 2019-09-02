package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type SyncConfig struct {
	system.Rsync
	TargetServer string
	TargetLabel  string
	//TargetServer           string // TODO this should not be mandatory
	//TargetLabel            string

}

func (s *SyncConfig) Init() error {
	if s.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if s.TargetServer == "" { // TODO remove
		return errs.With("TargetServer cannot be empty")
	}

	//targetPath := r.targetPath(*disks.FindDeepestBlockDeviceByLabel(r.TargetLabel))
	return nil
}
