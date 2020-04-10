package handler

import "github.com/n0rad/go-erlog/data"

type HandlerBuilder struct {
	//filter HandlerFilter
	New func() Handler
}

type Handler interface {
	Start() error
	Stop(err error)
	GetFields() data.Fields
	//HandleEvent(eventType EventType) error

	Add() error
	Remove() error
	Change() error
}

/////////////////////

type CommonHandler struct {
	HandlerName string
	fields      data.Fields
	StopChan    chan struct{}
}

func (h *CommonHandler) GetFields() data.Fields {
	return h.fields
}

func (h *CommonHandler) Init(fields data.Fields) {
	h.fields = fields
	h.StopChan = make(chan struct{})
}

func (h *CommonHandler) Start() error {
	<-h.StopChan
	return nil
}

func (h *CommonHandler) Stop(err error) {
	close(h.StopChan)
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
