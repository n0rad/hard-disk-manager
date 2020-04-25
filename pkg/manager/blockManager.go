package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
)

type BlockManager struct {
	CommonManager

	udev  *system.UdevService
	block system.BlockDevice
}

func (m *BlockManager) Init(parent Manager, block system.BlockDevice, udev *system.UdevService) {
	m.CommonManager.Init(parent, data.WithField("path", block.Path), &hdm.HDM)
	m.block = block
	m.udev = udev

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

	for _, child := range m.block.Children {
		manager := &BlockManager{}
		manager.Init(m, child, udev)
		m.children[child.Path] = manager
	}

	//// INIT children mananger for files if mountpoint
	//if block.Mountpoint != "" {
	//	f := &FsManager{}
	//	f.Init()
	//	m.children["fs"] = f
	//}

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
