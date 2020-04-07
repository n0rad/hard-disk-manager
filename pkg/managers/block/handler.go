package block

import (
	"github.com/n0rad/go-erlog/data"
)

// all available block handlers
var handlers = map[string]handler{}

type handler struct {
	filter HandlerFilter
	new    func() Handler
}

type EventType string

const (
	Add    EventType = "add"
	Remove EventType = "remove"
	Change EventType = "change"
)

type Handler interface {
	Init(manager *Manager)
	Start() error
	Stop(err error)
	GetFields() data.Fields
	//HandleEvent(eventType EventType) error

	Add() error
	Remove() error
	Change() error
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

type CommonHandler struct {
	handlerName string
	fields      data.Fields
	manager     *Manager
	stop        chan struct{}
}

func (h *CommonHandler) GetFields() data.Fields {
	return h.fields
}

func (h *CommonHandler) Init(manager *Manager) {
	h.fields = data.WithField("path", manager.block.Path).WithField("handler", h.handlerName)
	h.manager = manager
	h.stop = make(chan struct{})
}

func (h *CommonHandler) Start() error {
	<-h.stop
	return nil
}

func (h *CommonHandler) Stop(err error) {
	close(h.stop)
}

func (h *CommonHandler) Add() error {
	return nil
}

func (h *CommonHandler) Remove() error {
	return nil
}

func (h *CommonHandler) Change() error {
	return nil
}
