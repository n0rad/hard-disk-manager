package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
)

type BlockManager struct {
	CommonManager

	lsblk *system.Lsblk
	udev  *system.UdevService

	path  string
	block system.BlockDevice
}

func (m *BlockManager) Init(parent Manager, lsblk *system.Lsblk, path string, udev *system.UdevService) error {
	m.CommonManager.Init(parent, data.WithField("path", path), &hdm.HDM)
	m.lsblk = lsblk
	m.udev = udev
	m.path = path

	if err := m.updateBlockDeviceAndChildren(); err != nil {
		return err
	}

	// init handlers
	for name, builder := range blockHandlers {
		if builder.filter.Match(HandlerFilter{Type: m.block.Type, FSType: m.block.Fstype}) {
			handler := builder.new()
			handler.Init(name, m)
			logs.WithF(handler.GetFields()).Trace("new builder")

			m.handlers = append(m.handlers, handler)

			// TODO load configuration for builder
			// TODO if disabled, remove
		}
	}

	return nil
}

func (m *BlockManager) updateBlockDeviceAndChildren() error {
	block, err := m.lsblk.GetBlockDevice(m.path)
	if err != nil {
		return errs.WithE(err, "Failed to get block device to init manager")
	}
	m.block = block

	for _, child := range m.block.Children {
		found, ok := m.children[child.Path]
		if !ok {
			manager := &BlockManager{}
			if err := manager.Init(m, m.lsblk, child.Path, m.udev); err != nil {
				return err
			}
			m.children[child.Path] = manager
		} else {
			// TODO
			logs.WithF(m.fields.WithField("child-path", child.Path).WithField("found", found)).Fatal("Update already exists children")
		}
	}
	return nil
}

func (m *BlockManager) Start() error {
	udevChan := m.udev.Watch(m.block.Path)
	defer m.udev.Unwatch(udevChan)

	if err := m.preStart(); err != nil {
		return err
	}

	for {
		select {
		case event := <-udevChan:
			if event.Action == netlink.ADD {
				m.HandleEvent(Add)
			}
		case <-m.stop:
		}
	}

	return m.postStart()
}
