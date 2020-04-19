package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

type PathManager struct {
	CommonManager

	path string
	config PathConfig
}

func (m *PathManager) Init(path string) error {
	m.CommonManager.Init(data.WithField("path", path), &hdm.HDM)
	m.path = path
	if err := m.config.LoadFromDirIfExists(path); err != nil {
		return errs.WithE(err, "Failed to load config")
	}

	// init handlers
	for name, builder := range pathHandlerBuilders {
		handler := builder.new()
		handler.Init(name, m)

		m.handlers = append(m.handlers, handler)
	}

	return nil
}
