package manager

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type FsManager struct {
	PathManager

	block system.BlockDevice
	config   FsConfig
}

func (m *FsManager) Init(parent Manager, block system.BlockDevice) error {
	if err := m.PathManager.Init(parent, block.GetUsableLabel(), block.Mountpoint); err != nil {
		return errs.WithEF(err, m.fields, "Failed to init fs manager")
	}

	if err := m.config.LoadFromDirIfExists(block.Mountpoint); err != nil {
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


// TODO free
// TODO kill lsof
// TODO sync
