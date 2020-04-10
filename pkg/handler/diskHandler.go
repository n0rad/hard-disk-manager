package handler

// all available disk handler
var DiskHandlerBuilders = map[string]diskHandlerBuilder{}

type diskHandlerBuilder struct {
	new func() DiskHandler
}

type DiskHandler interface {
	Handler
	Init(name string, manager *DiskManager)
}


type CommonDiskHandler struct {
	CommonBlockHandler
}

func (h *CommonDiskHandler) Init(name string, manager *DiskManager) {
	h.CommonBlockHandler.Init(name, &manager.BlockManager)
}
