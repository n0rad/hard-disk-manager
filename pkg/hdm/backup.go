package hdm

import (
	"time"
)

const PathBackups = "/Backups"

type BackupConfig struct {
	TargetLabel string
	Interval    time.Duration
	Delete      bool
	//filter   string
}

type Backup struct {
	SourceDiskLabel string
	SourcePath      string
	Config          BackupConfig
	LastBackup      time.Time
}
