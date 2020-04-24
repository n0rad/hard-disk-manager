package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
)

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
	name     string
	fields   data.Fields
	stopChan chan struct{}
}

func (h *CommonHandler) GetFields() data.Fields {
	return h.fields
}

func (h *CommonHandler) Init(name string, fields data.Fields) {
	h.name = name
	h.fields = fields.WithField("name", h.name)
	h.stopChan = make(chan struct{})
	logs.WithF(h.GetFields()).Debug("New handler")
}

// called only once to start the handler
func (h *CommonHandler) Start() error {
	<-h.stopChan
	return nil
}

// This handler needs to be stopped
func (h *CommonHandler) Stop(err error) {
	close(h.stopChan)
}

func (h *CommonHandler) Add() error {
	return nil
}

// The handled layer needs to be removed
func (h *CommonHandler) Remove() error {
	return nil
}

func (h *CommonHandler) Change() error {
	return nil
}
