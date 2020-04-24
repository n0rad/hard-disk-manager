package handler

import "github.com/n0rad/go-erlog/data"

// all available block handler
var pathHandlerBuilders = map[string]pathHandlerBuilder{}

type pathHandlerBuilder struct {
	new func() PathHandler
}

type PathHandler interface {
	Handler
	Init(name string, manager *PathManager)
}

type CommonPathHandler struct {
	CommonHandler
	manager  *PathManager
}

func (h *CommonPathHandler) Init(name string, manager *PathManager) {
	h.CommonHandler.Init(name, data.WithField("path", manager.path))
	h.manager = manager
}
