package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type BlockManager struct {
	CommonManager

	Block system.BlockDevice
}

func (m *BlockManager) Init(block system.BlockDevice) {
	m.CommonManager.Init(data.WithField("path", block.Path), &hdm.HDM)
	m.Block = block

	// init builder
	for _, builder := range BlockHandlers {
		//if builder.filter.Match(HandlerFilter{Type: m.block.Type, FSType: m.block.Fstype}) {
		handler := builder.New()
		handler.Init(m)
		logs.WithF(handler.GetFields()).Trace("New builder")

		m.handlers = append(m.handlers, handler)

		// TODO load configuration for builder
		// TODO if disabled, remove
		//}
	}

	for _, child := range m.Block.Children {
		manager := &BlockManager{}
		manager.Init(child)
		m.children[child.Path] = manager
	}

	//// INIT children mananger for files if mountpoint
	//if block.Mountpoint != "" {
	//	f := &FileManager{}
	//	f.Init()
	//	m.children["fs"] = f
	//}

}