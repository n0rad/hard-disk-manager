package system

import (
	"github.com/alessio/shellescape"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
)

const hdmYamlFilename="hdm.yaml"

type HdmConfig struct {
	Backups []Backup

	configPath string
	fields     data.Fields
}

func (h *HdmConfig) Init(filesystem BlockDevice, configPath string) error {
	for i := range h.Backups {
		if err := h.Backups[i].Init(filesystem, configPath); err != nil {
			return err
		}
	}

	h.configPath = configPath
	h.fields = data.WithField("hdm", h.configPath)
	return nil
}

func (h *HdmConfig) FillFromFile(filesystem BlockDevice, file string) error {
	bytes, err := filesystem.server.Exec("sudo cat " + shellescape.Quote(file))
	if err != nil {
		return errs.WithEF(err, data.WithField("file", file),"Failed to cat file")
	}

	if err := yaml.Unmarshal(bytes, h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to parse hdm file")
	}

	if err := h.Init(filesystem, file); err != nil {
		return errs.WithEF(err, h.fields.WithField("content", string(bytes)), "Failed to init hdm file")
	}

	return nil
}

func (h *HdmConfig) RunBackups(disks Disks) error {
	for _, backup := range h.Backups {
		if err := backup.Backup(disks); err != nil {
			return err
		}
	}
	return nil
}