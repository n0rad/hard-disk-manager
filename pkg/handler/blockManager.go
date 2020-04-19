package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type BlockManager struct {
	CommonManager

	block system.BlockDevice
}

func (m *BlockManager) Init(block system.BlockDevice) {
	m.CommonManager.Init(data.WithField("path", block.Path), &hdm.HDM)
	m.block = block

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
		manager.Init(child)
		m.children[child.Path] = manager
	}

	//// INIT children mananger for files if mountpoint
	//if block.Mountpoint != "" {
	//	f := &FsManager{}
	//	f.Init()
	//	m.children["fs"] = f
	//}

}
