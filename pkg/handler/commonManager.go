package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

type EventType string

const (
	Add    EventType = "add"
	Remove EventType = "remove"
	Change EventType = "change"
)

type Manager interface {
	GetHDM() *hdm.Hdm
	HandleEvent(eventType EventType)
	//Parent() Manager
	Start() error
	Stop(err error)
}

type CommonManager struct {
	hdm    *hdm.Hdm
	fields data.Fields
	handlers []Handler
	children map[string]Manager
	stop     chan struct{}
}

func (m *CommonManager) GetHDM() *hdm.Hdm {
	return m.hdm
}

func (m *CommonManager) Init(fields data.Fields, hdm *hdm.Hdm) {
	logs.WithF(fields).Info("new manager")
	m.children = map[string]Manager{}
	m.hdm = hdm
	m.fields = fields
}

func (m *CommonManager) HandleEvent(eventType EventType) {
	switch eventType {
	case Add:
		// going downstream
		for _, h := range m.handlers {
			if err := h.Add(); err != nil {
				logs.WithEF(err, data.WithField("event", eventType)).Error("Failed to handle event")
			}
		}
		for _, subManager := range m.children {
			subManager.HandleEvent(eventType)
		}
	case Change:
		// going downstream
		// TODO
	case Remove:
		// going upstream
		for _, subManager := range m.children {
			subManager.HandleEvent(eventType)
		}
		for _, h := range m.handlers {
			if err := h.Remove(); err != nil {
				logs.WithEF(err, data.WithField("event", eventType)).Error("Failed to handle event")
			}
		}
	}
}

func (m *CommonManager) Start() error {
	m.stop = make(chan struct{})

	for _, h := range m.handlers {
		logs.WithF(h.GetFields()).Trace("Starting handler")
		go h.Start()
	}

	for c := range m.children {
		manager := m.children[c]
		go manager.Start()
	}

	<-m.stop

	for _, h := range m.handlers {
		h.Stop(nil)
	}
	return nil
}

func (m *CommonManager) Stop(error) {
	close(m.stop)
}
