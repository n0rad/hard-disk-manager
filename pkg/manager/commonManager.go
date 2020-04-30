package manager

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

type EventType string

const (
	Add    EventType = "add"
	Remove EventType = "remove"
	Change EventType = "change" // new children
)

type Manager interface {
	HandleEvent(eventType EventType)
	Start() error
	Stop(err error)

	runSerialJob(f func() interface{}) <-chan interface{} // reduce pressure on disk
}

type CommonManager struct {
	hdm        *hdm.Hdm
	fields     data.Fields
	handlers   []Handler
	parent     Manager
	children   map[string]Manager
	stop       chan struct{}
}

func (m *CommonManager) Init(parent Manager, fields data.Fields, hdm *hdm.Hdm) {
	logs.WithF(fields).Debug("new manager")
	m.children = map[string]Manager{}
	m.hdm = hdm
	m.fields = fields
	m.parent = parent
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
		for k, subManager := range m.children {
			subManager.HandleEvent(eventType)
			delete(m.children, k)
		}
		for _, h := range m.handlers {
			if err := h.Remove(); err != nil {
				logs.WithEF(err, data.WithField("event", eventType)).Error("Failed to handle event")
			}
		}
		m.Stop(nil)
	}
}

func (m *CommonManager) preStart() error {
	m.stop = make(chan struct{})

	for _, h := range m.handlers {
		logs.WithF(h.GetFields()).Trace("Starting handler")
		go h.Start()
	}

	for c := range m.children {
		manager := m.children[c]
		go manager.Start()
	}
	return nil
}

func (m *CommonManager) postStart() error {
	for s := range m.children {
		m.children[s].Stop(nil)
		delete(m.children, s)
	}

	for _, h := range m.handlers {
		h.Stop(nil)
	}
	return nil
}

func (m *CommonManager) Start() error {
	if err := m.preStart(); err != nil {
		return err
	}

	<-m.stop

	return m.postStart()
}

func (m *CommonManager) Stop(e error) {
	logs.WithF(m.fields).Trace("Stopping")
	close(m.stop)
}

/////////////////////////////////////////////////////

func (m *CommonManager) runSerialJob(f func() interface{}) <-chan interface{} {
	return m.parent.runSerialJob(f)
}
