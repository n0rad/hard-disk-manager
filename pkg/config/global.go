package config

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const pathConfigYaml = "/config.yaml"

type GlobalConfig struct {
	Defaults PathConfig
}

func (g *GlobalConfig) Init(hdmHome string) error {

	configPath := hdmHome + pathConfigYaml

	if _, err := os.Stat(configPath); err == nil {
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errs.WithEF(err, data.WithField("file", configPath), "Failed to read configuration file")
		}

		if err = yaml.Unmarshal(file, g); err != nil {
			return errs.WithEF(err, data.WithField("file", configPath), "Invalid configuration format")
		}
	}

	if err := g.Defaults.Init(); err != nil {
		return errs.WithE(err, "Failed to init default config")
	}

	return nil
}

func (g GlobalConfig) GetPathConfig(path string) (PathConfig, error) {
	// TODO deepcopy defaults
	conf := g.Defaults

	absPath, err := filepath.Abs(path)
	if err != nil {
		return conf, errs.WithEF(err, data.WithField("path", path), "Failed to convert to absolute path")
	}

	file := findHdmYamlInParentFolder(absPath)
	if file != "" {
		if err := conf.Load(file); err != nil {
			return conf, errs.WithEF(err, data.WithField("file", file), "Failed to load path config")
		}
	}
	return conf, nil
}

////////////////////////////

func findHdmYamlInParentFolder(path string) string {
	current := path
	for current != "/" {
		fullPath := current + PathHdmYaml
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
		current = filepath.Dir(current)
	}
	return ""
}
