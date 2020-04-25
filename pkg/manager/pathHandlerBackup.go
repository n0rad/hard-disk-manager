package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

func init() {
	pathHandlerBuilders["backup"] = pathHandlerBuilder{
		new: func() PathHandler {
			return &PathHandlerBackup{}
		},
	}
}

type PathHandlerBackup struct {
	CommonPathHandler
}

func (h *PathHandlerBackup) Add() error {
	for _, v := range h.manager.config.Backups {
		go h.handleBackup(v)
	}
	return nil
}

func (h *PathHandlerBackup) handleBackup(config hdm.BackupConfig) {
	for {
		backup, err := h.manager.hdm.BackupDB().GetOrCreateBackup(h.manager.diskLabel, h.manager.path)
		if err != nil {
			logs.WithEF(err, h.fields).Error("Failed to get backup from db")
			return
		}

		backup.Config = config
		runner := BackupRunner{
			backup: backup,
			rsync:  system.Rsync{},
		}

		select {
		case <-h.stop:
			return
		case <-time.After(backup.LastBackup.Add(config.Interval).Sub(time.Now())):
			res := <-h.manager.runSerialJob(func() interface{} {
				logs.WithFields(h.fields).Info("Time to backup")
				if err := runner.Backup(); err != nil {
					return errs.WithEF(err, h.fields, "Backup failed")
				}
				return nil
			})
			if err, _ :=res.(error); err != nil {
				logs.WithEF(err, h.fields).Error("Backup failed, retry in 10min")
				select {
				case <-h.stop:
					return
				case <-time.After(10*time.Minute):
				}
			}
		}
	}
}


////////////////////////

type BackupRunner struct {
	backup hdm.Backup
	rsync system.Rsync
}

//func (b *BackupRunner) Init(configPath string, sourceBlockDevice system.BlockDevice, servers hdm.Servers) error {
func (b *BackupRunner) Init(backup hdm.Backup, servers hdm.Servers) error {
	if b.backup.Config.TargetLabel == "" {
		return errs.With("TargetLabel cannot be empty")
	}

	//// source
	//sourcePath, err := securejoin.SecureJoin(filepath.Dir(configPath), b.Path)
	//if err != nil {
	//	return errs.WithEF(err, data.WithField("path", b.Path), "Failed to prepare path")
	//}

	// target
	targetBlockDevice, err := servers.GetBlockDeviceByLabel(b.backup.Config.TargetLabel)
	if err != nil {
		return errs.WithE(err, "Cannot start backup, target device not found")
	}
	if targetBlockDevice.Mountpoint == "" {
		return errs.WithF(data.WithField("blockDevice", targetBlockDevice), "target filesystem is not mounted")
	}
	//targetPath, err := securejoin.SecureJoin(targetBlockDevice.Mountpoint+hdm.PathBackups, b.Path)
	//if err != nil {
	//	return errs.WithEF(err, data.WithField("path", targetBlockDevice.Mountpoint+hdm.PathBackups+b.Path), "Failed to prepare tatget path")
	//}

	//b.fields = data.WithField("hdm", b.configPath).WithField("TargetLabel", b.TargetLabel)

	//b.rsync = system.Rsync{
	//	SourceInFilesystemPath: backup.SourcePath,
	//	SourceFilesystem:       sourceBlockDevice,
	//	TargetInFilesystemPath: targetBlockDevice.Mountpoint+hdm.PathBackups+"/dddddd",
	//	TargetFilesystem:       targetBlockDevice,
	//	Delete:                 b.backup.Config.Delete,
	//}
	//
	//if err := b.rsync.Init(); err != nil {
	//	return errs.WithE(err, "Failed to init rsync")
	//}
	return nil
}

func (b *BackupRunner) Backupable() (error, error) {
	return b.rsync.Rsyncable()
}

func (b *BackupRunner) Backup() error {
	return b.rsync.RSync()
}
