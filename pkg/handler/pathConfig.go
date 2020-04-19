package handler

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"io/ioutil"
	"os"
)

const HdmYamlFilename = "hdm.yaml"
const PathHdmYaml = "/" + HdmYamlFilename

type PathConfig struct {
	Backups []hdm.BackupConfig
	//Syncs   []hdm.SyncConfig

	configPath string
}

func (h *PathConfig) LoadFromDirIfExists(directory string) error {
	return loadFromDirIfExistsToStruct(directory, h)
}

func loadFromDirIfExistsToStruct(directory string, cfg interface{}) error {
	configPath := directory+PathHdmYaml

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errs.WithEF(err, data.WithField("path", configPath),"Cannot identify if config file exists")
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", configPath), "Failed to read hdm config file")
	}

	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)).WithField("path", configPath), "Failed to parse hdm file")
	}
	return nil
}
