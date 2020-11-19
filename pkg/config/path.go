package config

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"io/ioutil"
)

const HdmYaml = "hdm.yaml"
const PathHdmYaml = "/" + HdmYaml

type PathConfig struct {
	Checksum ChecksumConfig
}

func (h *PathConfig) Init() error {
	if err := h.Checksum.Init(); err != nil {
		return errs.WithE(err, "Failed to init checksum config")
	}
	return nil
}

func (h *PathConfig) Load(configPath string) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", configPath), "Failed to read hdm config file")
	}

	if err := yaml.Unmarshal(bytes, h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)).WithField("path", configPath), "Failed to parse hdm file")
	}

	if err := h.Init(); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to init hdm file")
	}
	return nil
}

//func (h *PathConfig) GetConfigForPath() (PathConfig, error) {
//
//}
//
//func FindConfigInParentFolder(path string) (PathConfig, error) {
//	//conf := findHdmYamlInParentFolder(path)
//	//if conf != nil {
//	//	config := PathConfig{}
//	//	config.Load(path)
//	//	config.Init(path)
//	//}
//	//return nil
//}
