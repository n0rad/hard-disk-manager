package handler

// all available disk handler
var DiskHandlerBuilders = map[string]diskHandlerBuilder{}

type diskHandlerBuilder struct {
	New func() DiskHandler
}

type DiskHandler interface {
	Handler
	Init(manager *DiskManager)
}


type CommonDiskHandler struct {
	CommonBlockHandler
}

func (h *CommonDiskHandler) Init(manager *DiskManager) {
	h.CommonBlockHandler.Init(&manager.BlockManager)
}
