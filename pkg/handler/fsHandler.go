package handler

// all available file handler
var fsHandlerBuilders = map[string]fsHandlerBuilder{}

type fsHandlerBuilder struct {
	new  func() FsHandler
}

type FsHandler interface {
	Handler
	Init(name string, manager *FsManager)
}

type CommonFsHandler struct {
	CommonPathHandler
	manager *FsManager
}

func (h *CommonFsHandler) Init(name string, manager *FsManager) {
	h.manager = manager
	h.CommonPathHandler.Init(name, &manager.PathManager)
}
