package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

const pathBackups = "/Backups"

type BackupConfig struct {
	TargetLabel string
	//Interval time.Duration
	//Filter   string
	Path   string
	Delete bool

	//configPath string
	rsync system.Rsync
}

type Backup struct {
	Config     BackupConfig
	Path       string
	DiskName   string
	LastBackup time.Time
}

func (b *BackupConfig) Init(configPath string) error {
	b.rsync = system.Rsync{
		SourceInFilesystemPath: "",
		//SourceFilesystem: ,
		TargetInFilesystemPath: "",
		//TargetFilesystem: "",
		Delete: b.Delete,
	}
	if b.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}
	if len(b.Path) > 0 && b.Path[0] != '/' {
		b.Path = "/" + b.Path
	}
	//b.configPath = configPath
	//b.fullPath = filepath.Dir(configPath) + b.Path // TODO remove ../
	//b.inBlockDevicePath = strings.Replace(b.fullPath, filesystem.Mountpoint, "", 1)
	//b.filesystem = filesystem
	//b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)
	return nil
}

func (b *BackupConfig) Backupable(server system.Server) (error, error) {
	return nil, nil
}

func (b *BackupConfig) Backup(server system.Server) error {
	return nil
}
