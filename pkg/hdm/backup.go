package hdm

import (
	"github.com/cyphar/filepath-securejoin"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"path/filepath"
)

const PathBackups = "/Backups"

type BackupConfig struct {
	TargetLabel string
	//Interval time.Duration
	//Filter   string
	Path   string
	Delete bool
}

type Backup struct {
	BackupConfig
	//Path       string
	//DiskName   string
	//LastBackup time.Time

	//configPath string
	rsync system.Rsync
}

func (b *Backup) Init(configPath string, sourceBlockDevice system.BlockDevice, servers Servers) error {

	if b.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if len(b.Path) > 0 && b.Path[0] != '/' {
		b.Path = "/" + b.Path
	}

	// source
	sourcePath, err := securejoin.SecureJoin(filepath.Dir(configPath), b.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", b.Path), "Failed to prepare path")
	}

	// target
	targetBlockDevice, err := servers.GetBlockDeviceByLabel(b.TargetLabel)
	if err != nil {
		return errs.WithE(err, "Cannot start backup, target device not found")
	}
	if targetBlockDevice.Mountpoint == "" {
		return errs.WithF(data.WithField("blockDevice", targetBlockDevice), "target filesystem is not mounted")
	}
	targetPath, err := securejoin.SecureJoin(targetBlockDevice.Mountpoint+PathBackups, b.Path)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", targetBlockDevice.Mountpoint+PathBackups+b.Path), "Failed to prepare tatget path")
	}

	//b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)

	b.rsync = system.Rsync{
		SourceInFilesystemPath: sourcePath,
		SourceFilesystem:       sourceBlockDevice,
		TargetInFilesystemPath: targetPath,
		TargetFilesystem:       targetBlockDevice,
		Delete:                 b.Delete,
	}

	if err := b.rsync.Init(); err != nil {
		return errs.WithE(err, "Failed to init rsync")
	}
	return nil
}

func (b *Backup) Backupable() (error, error) {
	return b.rsync.Rsyncable()
}

func (b *Backup) Backup() error {
	return b.rsync.RSync()
}
