package hdm

import (
	"github.com/alessio/shellescape"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"os"
	"strings"
)

const HdmYamlFilename = "hdm.yaml"
const PathHdmYaml = "/" + HdmYamlFilename

type Config struct {
	Backups []BackupConfig
	Syncs   []SyncConfig

	RecursiveConfig bool

	configPath string
	fields     data.Fields
}

func NewConfig(hdmConfigPath string) (Config, error) {
	var cfg Config

	bytes, err := ioutil.ReadFile(hdmConfigPath)
	if err != nil {
		return cfg, errs.WithEF(err, data.WithField("path", hdmConfigPath), "Failed to read hdm config file")
	}

	if err := yaml.Unmarshal([]byte(bytes), &cfg); err != nil {
		return cfg, errs.WithEF(err, data.WithField("content", string(bytes)).WithField("path", hdmConfigPath), "Failed to parse hdm file")
	}

	if err := cfg.Init(hdmConfigPath); err != nil {
		return cfg, errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to init hdm file")
	}

	return cfg, nil
}

func (h Config) GetConfigPath() string {
	return h.configPath
}

func (h *Config) Init(configPath string) error {
	for i := range h.Backups {
		if err := h.Backups[i].Init(configPath); err != nil {
			return err
		}
	}

	h.configPath = configPath
	h.fields = data.WithField("hdm", h.configPath)
	return nil
}

func (h *Config) FillFromFile(filesystem system.BlockDeviceOLD, file string) error {
	bytes, err := filesystem.ExecShell("cat " + shellescape.Quote(file))
	if err != nil {
		return errs.WithEF(err, data.WithField("file", file), "Failed to cat file")
	}

	if err := yaml.Unmarshal([]byte(bytes), h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to parse hdm file")
	}

	if err := h.Init(file); err != nil {
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

//if len(b.Children) > 0 {
//	for _, child := range b.Children {
//		configs, err := hdm.FindConfigs(child)
//		if err != nil {
//			return hdmConfigs, err
//		}
//		hdmConfigs = append(hdmConfigs, configs...)
//	}
//	return hdmConfigs, nil
//}

func FindConfigs(path string, server system.Server) ([]Config, error) {
	var hdmConfigs []Config

	if path == "" {
		return hdmConfigs, errs.WithF(data.WithField("path", path), "BlockDeviceOLD is not mounted")
	}

	hdmRootFilePath := path + PathHdmYaml

	if _, err := os.Stat(hdmRootFilePath); err != nil {
		logs.WithEF(err, data.WithField("path", path)).Debug("hdm root file does not exists or cannot be read")
		return hdmConfigs, nil
	}

	config, err := NewConfig(hdmRootFilePath)
	if err != nil {
		logs.WithEF(err, data.WithField("path", hdmRootFilePath)).Debug("Failed to read hdm root file")
		return hdmConfigs, nil
	}

	hdmConfigs = append(hdmConfigs, config)
	if !config.RecursiveConfig {
		return hdmConfigs, nil
	}

	configs, err := server.Exec("find", path, "-type", "f", "-not", "-path", path+pathBackups+"/*", "-name", HdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithE(err, "Failed to find hdm.yaml files")
	}

	lines := strings.Split(string(configs), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		logs.WithField("path", line).Debug(HdmYamlFilename + " found")
		if cfg, err := NewConfig(line); err != nil {
			logs.WithField("path", line).Warn("Failed to read hdm.yaml configuration")
		} else {
			hdmConfigs = append(hdmConfigs, cfg)
		}
	}
	return hdmConfigs, nil
}
