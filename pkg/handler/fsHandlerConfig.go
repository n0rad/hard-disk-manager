package handler

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"path"
	"strings"
)

func init() {
	fsHandlerBuilders["config"] = fsHandlerBuilder{
		new: func() FsHandler {
			return &FsHandlerFiles{}
		},
	}
}

type FsHandlerConfig struct {
	CommonFsHandler
}

func (h *FsHandlerConfig) Add() error {
	if h.manager.config.SearchConfigs {
		configs, err := h.findConfigs(h.manager.path, runner.Local)
		if err != nil {
			return errs.WithEF(err, h.fields, "Failed to find configs")
		}

		for _, configPath := range configs {
			configDir := path.Dir(configPath)

			if configDir == h.manager.path {
				continue
			}

			m := PathManager{}
			if err := m.Init(&h.manager.CommonManager, h.manager.diskLabel, configDir); err != nil {
				return err
			}
			h.manager.children[configDir] = &m
		}

	}

	return nil
}

func (h FsHandlerConfig) findConfigs(rootPath string, exec runner.Exec) ([]string, error) {
	var hdmConfigs []string

	configs, err := exec.ExecGetStdout("find", rootPath, "-type", "f", "-not", "-path", rootPath+hdm.PathBackups+"/*", "-name", HdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithE(err, "Failed to find "+HdmYamlFilename+" files")
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
