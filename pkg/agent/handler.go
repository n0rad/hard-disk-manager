package agent

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type Handler interface {
	Init(manager *DiskManager)
	Start()
	Stop()
	//Handle(event)
}

type CommonHandler struct {
	disk    system.Disk
	server  system.Server
	fields  data.Fields
	manager *DiskManager
}

func (h *CommonHandler) Init(manager *DiskManager) {
	h.fields = data.WithField("path", manager.Path)
	h.server = system.Server{}
	h.manager = manager

	if err := h.server.Init(); err != nil {
		logs.WithE(err).Error("fail")
	}
}

func (h *CommonHandler) Stop() {

}
