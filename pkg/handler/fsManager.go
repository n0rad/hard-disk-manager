package handler

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
)

type FsManager struct {
	PathManager

	config   FsConfig
}

func (m *FsManager) Init(rootPath string, exec runner.Exec) error {
	// TODO use block device instead ???
	m.PathManager.Init(rootPath)

	m.path = rootPath
	if err := m.config.LoadFromDirIfExists(rootPath); err != nil {
		return errs.WithE(err, "Failed to load config")
	}

	// init handlers
	for name, builder := range fsHandlerBuilders {
		handler := builder.new()
		handler.Init(name, m)

		m.handlers = append(m.handlers, handler)
	}

	return nil
}
