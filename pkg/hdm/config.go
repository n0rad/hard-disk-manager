package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"io/ioutil"
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

func (h Config) GetConfigPath() string {
	return h.configPath
}

func (h *Config) Init(configPath string) error {
	h.configPath = configPath
	h.fields = data.WithField("hdm", h.configPath)
	return nil
}

func (h *Config) Load(configPath string) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", configPath), "Failed to read hdm config file")
	}

	if err := yaml.Unmarshal(bytes, h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)).WithField("path", configPath), "Failed to parse hdm file")
	}

	if err := h.Init(configPath); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to init hdm file")
	}
	return nil
}

/////////////////


// TODO remove
func NewConfig(hdmConfigPath string) (Config, error) {
	var cfg Config
	return cfg, cfg.Load(hdmConfigPath)
}

func FindConfigs(rootPath string, exec runner.Exec) ([]string, error) {
	var hdmConfigs []string

	configs, err := exec.ExecGetStdout("find", rootPath, "-type", "f", "-not", "-path", rootPath+PathBackups+"/*", "-name", HdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithE(err, "Failed to find hdm.yaml files")
	}

	lines := strings.Split(configs, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		hdmConfigs = append(hdmConfigs, line)
	}
	return hdmConfigs, nil
}
