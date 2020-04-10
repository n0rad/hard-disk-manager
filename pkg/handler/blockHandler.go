package handler

import (
	"github.com/n0rad/go-erlog/data"
)

// all available Block handler
var BlockHandlers = map[string]BlockHandlerBuilder{}

type BlockHandlerBuilder struct {
	Filter HandlerFilter
	New func() BlockHandler
}

type BlockHandler interface {
	Handler
	Init(manager *BlockManager)
}

//////////////////

type CommonBlockHandler struct {
	CommonHandler
	manager *BlockManager
}

func (h *CommonBlockHandler) Init(manager *BlockManager) {
	h.CommonHandler.Init(data.WithField("path", manager.Block.Path).WithField("handler", h.HandlerName))
	h.manager = manager
}

/////////////////

type HandlerFilter struct {
	Type   string
	FSType string
}

func (h HandlerFilter) Match(filter HandlerFilter) bool {
	typeRes := true
	if h.Type != "" {
		if filter.Type == h.Type {
			typeRes = true
		} else {
			typeRes = false
		}
	}

	fsTypeRes := true
	if h.FSType != "" {
		if filter.FSType == h.FSType {
			fsTypeRes = true
		} else {
			fsTypeRes = false
		}
	}
	return typeRes && fsTypeRes
}

///////////////////////////
