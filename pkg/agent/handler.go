package agent

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type Handler interface {
	Init(path string)
	Start()
	Stop()
	//Handle(event)
}

type CommonHandler struct {
	disk system.Disk

	path   string
	server system.Server
	fields data.Fields
}

func (h *CommonHandler) Init(path string) {
	h.path = path
	h.fields = data.WithField("path", path)
	h.server = system.Server{}

	if err := h.server.Init(); err != nil {
		logs.WithE(err).Error("fail")
	}
}

func (h *CommonHandler) Stop() {

}
