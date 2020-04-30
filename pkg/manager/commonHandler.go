package manager

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
	name   string
	fields data.Fields
	stop   chan struct{}
}

func (h *CommonHandler) GetFields() data.Fields {
	return h.fields
}

func (h *CommonHandler) Init(name string, fields data.Fields) {
	h.name = name
	h.fields = fields.WithField("name", h.name)
	h.stop = make(chan struct{})
	logs.WithF(h.GetFields()).Debug("New manager")
}

// called only once to start the manager
func (h *CommonHandler) Start() error {
	<-h.stop
	return nil
}

// This manager needs to be stopped
func (h *CommonHandler) Stop(err error) {
	logs.WithF(h.fields).Debug("Closing")
	close(h.stop)
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
