package agent

import (
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
	path string
	server system.Server
}

func (h *CommonHandler) Init(path string) {
	h.path = path
	h.server = system.Server{
	}
	if err := h.server.Init(); err != nil {
		logs.WithE(err).Error("fail")
	}
}



type Event struct {

}
