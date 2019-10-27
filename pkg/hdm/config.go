package hdm

import (
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
	h.configPath = configPath
	h.fields = data.WithField("hdm", h.configPath)
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

func FindConfigs(blockdevice system.BlockDevice, server Server) ([]Config, error) {
	var hdmConfigs []Config

	if blockdevice.Mountpoint == "" {
		return hdmConfigs, errs.WithF(data.WithField("path", blockdevice.Mountpoint), "BlockDevice is not mounted")
	}

	hdmRootFilePath := blockdevice.Mountpoint + PathHdmYaml

	if _, err := os.Stat(hdmRootFilePath); err != nil {
		logs.WithEF(err, data.WithField("path", blockdevice.Mountpoint)).Debug("hdm root file does not exists or cannot be read")
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

	configs, err := blockdevice.GetExec().ExecGetStdout("find", blockdevice.Mountpoint, "-type", "f", "-not", "-path", blockdevice.Mountpoint+PathBackups+"/*", "-name", HdmYamlFilename)
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
