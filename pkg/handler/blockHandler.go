package handler

import (
	"github.com/n0rad/go-erlog/data"
)

// all available Block handler
var BlockHandlers = map[string]BlockHandlerBuilder{}

type BlockHandlerBuilder struct {
	Name string
	Filter HandlerFilter
	New func() BlockHandler
}

type BlockHandler interface {
	Handler
	Init(name string, manager *BlockManager)
}

//////////////////

type CommonBlockHandler struct {
	CommonHandler
	manager *BlockManager
}

func (h *CommonBlockHandler) Init(name string, manager *BlockManager) {
	h.CommonHandler.Init(name, data.WithField("path", manager.Block.Path))
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
