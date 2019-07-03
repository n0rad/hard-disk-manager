package hdm

import (
	"github.com/alessio/shellescape"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const pathBackups = "/Backups"

type BackupConfig struct {
	TargetLabel string
	//Interval time.Duration
	//Filter   string
	Path   string
	Delete bool

	fields            data.Fields
	fullPath          string
	inBlockDevicePath string
	configPath        string
	filesystem        system.BlockDevice
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
	b.fullPath = filepath.Dir(configPath) + b.Path // TODO remove ../
	b.inBlockDevicePath = strings.Replace(b.fullPath, filesystem.Mountpoint, "", 1)
	b.filesystem = filesystem
	b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)
	return nil
}

func (b *BackupConfig) sourceSize() (int, error) {
	bytes, err := b.filesystem.Exec("sudo du -s " + shellescape.Quote(b.fullPath) + " | cut -f1")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to get directory size")
	}
	size, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to parse 'du' result")
	}
	return size, nil
}

func (b *BackupConfig) targetSize(target system.BlockDevice) (int, error) {
	targetPath := b.targetPath(target) + "/" + path.Base(b.fullPath)
	_, err := b.filesystem.Exec("sudo test -d " + shellescape.Quote(targetPath))
	if err != nil {
		return 0, nil
	}

	bytes, err := b.filesystem.Exec("sudo du -s " + shellescape.Quote(targetPath) + " | cut -f1")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to get directory size")
	}
	size, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to parse 'du' result")
	}
	return size, nil
}

func (b *BackupConfig) Backupable(disks system.Disks) (error, error) {
	target := disks.FindDeepestBlockDeviceByLabel(b.TargetLabel)

	if target == nil {
		return errs.WithF(b.fields, "Disk cannot be found"), nil
	}
	if target.Mountpoint == "" {
		return errs.WithF(b.fields.WithField("disk", target), "Disk is not mounted"), nil
	}

	sourceSize, err := b.sourceSize()
	if err != nil {
		return nil, errs.WithEF(err, b.fields, "Cannot get directory size")
	}

	targetSize, err := b.targetSize(*target)
	if err != nil {
		return nil, errs.WithEF(err, b.fields, "Cannot get directory size")
	}

	targetAvailable, err := target.SpaceAvailable()
	if err != nil {
		return nil, errs.WithEF(err, b.fields, "Cannot get TargetLabel available space")
	}

	if sourceSize > targetSize+targetAvailable {
		return errs.WithF(data.WithField("sourceSize", sourceSize).
			WithField("targetSize", targetSize).
			WithField("targetAvailable", targetAvailable), "Not enough space to backup"), nil
	}
	return nil, nil
}

func (b *BackupConfig) Backup(disks system.Disks) error {
	why, err := b.Backupable(disks)
	if err != nil {
		return errs.WithEF(err, b.fields, "Failed to see if directory is backupable")
	}
	if why != nil {
		logs.WithEF(why, b.fields).Warn("Directory is not backupable")
		return nil
	}

	targetPath := b.targetPath(*disks.FindDeepestBlockDeviceByLabel(b.TargetLabel))

	if _, err := b.filesystem.Exec("sudo mkdir -p " + targetPath); err != nil {
		return errs.WithEF(err, b.fields.WithField("path", targetPath), "Failed to create target backup path")
	}

	deleteIfSourceRemoved := ""
	if b.Delete {
		deleteIfSourceRemoved = "--delete"
	}

	logs.WithField("path", b.fullPath).WithField("target", b.TargetLabel).Info("Running backup")
	_, err = b.filesystem.Exec("sudo rsync -avP " + deleteIfSourceRemoved + " --itemize-changes " + shellescape.Quote(b.fullPath) + " " + targetPath) // TODO support sync to other server
	if err != nil {
		return errs.WithEF(err, b.fields, "Backup failed")
	}
	return nil
}

func (b *BackupConfig) targetPath(target system.BlockDevice) string {
	return shellescape.Quote(target.Mountpoint + pathBackups + path.Dir(b.inBlockDevicePath))
}
