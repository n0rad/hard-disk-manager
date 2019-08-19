package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	system "github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

const pathBackups = "/Backups"

type BackupConfig struct {
	TargetLabel string
	//Interval time.Duration
	//Filter   string
	Path   string
	Delete bool

	configPath        string
}

type Backup struct {
	Config     BackupConfig
	Path       string
	DiskName   string
	LastBackup time.Time
}

func (b *BackupConfig) Init(filesystem system.BlockDevice, configPath string) error {
	if b.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if len(b.Path) > 0 && b.Path[0] != '/' {
		b.Path = "/" + b.Path
	}
	b.configPath = configPath
	//b.fullPath = filepath.Dir(configPath) + b.Path // TODO remove ../
	//b.inBlockDevicePath = strings.Replace(b.fullPath, filesystem.Mountpoint, "", 1)
	//b.filesystem = filesystem
	//b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)
	return nil
}

func (b *BackupConfig) Backupable(disks system.Disks) (error, error) {
	return nil, nil
}

func (b *BackupConfig) Backup(disks system.Disks) error {
	return nil
}


