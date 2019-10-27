package hdm

import (
	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"path/filepath"
)

type SyncConfig struct {
	TargetServer string
	TargetLabel  string
	Path         string
}

func (s *SyncConfig) Init() error {
	if s.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if s.TargetServer == "" {
		return errs.With("TargetServer cannot be empty")
	}

	//targetPath := r.targetPath(*disks.FindDeepestBlockDeviceByLabel(r.TargetLabel))
	return nil
}

type Sync struct {
	SyncConfig

	rsync system.Rsync
}

func (s *Sync) Init(configPath string, sourceBlockDevice system.BlockDevice, servers Servers) error {

	if s.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if len(s.Path) > 0 && s.Path[0] != '/' {
		s.Path = "/" + s.Path
	}

	// source
	sourcePath, err := securejoin.SecureJoin(filepath.Dir(configPath), s.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", s.Path), "Failed to prepare path")
	}

	// target
	//TODO find target

	targetBlockDevice, err := servers.GetBlockDeviceByLabel(s.TargetLabel)
	if err != nil {
		return errs.WithE(err, "Cannot start sync, target device not found")
	}
	if targetBlockDevice.Mountpoint == "" {
		return errs.WithF(data.WithField("blockDevice", targetBlockDevice), "target filesystem is not mounted")
	}
	targetPath, err := securejoin.SecureJoin(targetBlockDevice.Mountpoint+PathBackups, s.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", targetBlockDevice.Mountpoint+PathBackups+s.Path), "Failed to prepare tatget path")
	}

	//b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)

	s.rsync = system.Rsync{
		SourceInFilesystemPath: sourcePath,
		SourceFilesystem:       sourceBlockDevice,
		TargetInFilesystemPath: targetPath,
		TargetFilesystem:       targetBlockDevice,
		Delete:                 true,
	}

	if err := s.rsync.Init(); err != nil {
		return errs.WithE(err, "Failed to init rsync")
	}
	return nil
}

func (s *Sync) Backupable() (error, error) {
	return s.rsync.Rsyncable()
}

func (s *Sync) Backup() error {
	return s.rsync.RSync()
}
