package hdm

import (
	"github.com/alessio/shellescape"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

const hdmYamlFilename="hdm.yaml"

type Config struct {
	Backups []BackupConfig
	Syncs []SyncConfig

	configPath string
	fields     data.Fields
}

func (h *Config) Init(filesystem system.BlockDevice, configPath string) error {
	for i := range h.Backups {
		if err := h.Backups[i].Init(filesystem, configPath); err != nil {
			return err
		}
	}

	h.configPath = configPath
	h.fields = data.WithField("hdm", h.configPath)
	return nil
}

func (h *Config) FillFromFile(filesystem system.BlockDevice, file string) error {
	bytes, err := filesystem.ExecShell("cat " + shellescape.Quote(file))
	if err != nil {
		return errs.WithEF(err, data.WithField("file", file),"Failed to cat file")
	}

	if err := yaml.Unmarshal([]byte(bytes), h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to parse hdm file")
	}

	if err := h.Init(filesystem, file); err != nil {
		return errs.WithEF(err, h.fields.WithField("content", string(bytes)), "Failed to init hdm file")
	}

	return nil
}

func (h *Config) RunBackups(disks system.Disks) error {
	for _, backup := range h.Backups {
		if err := backup.Backup(disks); err != nil {
			return err
		}
	}
	return nil
}