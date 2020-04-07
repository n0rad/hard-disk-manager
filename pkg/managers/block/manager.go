package block

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type Manager struct {
	hdm      *hdm.Hdm
	block    system.BlockDevice
	handlers []Handler
	children map[string]Manager
	stop       chan struct{}
}

func (m *Manager) Init(block system.BlockDevice, hdm *hdm.Hdm) {
	logs.WithField("path", block.Path).Info("New block manager")
	m.children = map[string]Manager{}
	m.block = block
	m.hdm = hdm

	// init handlers
	for _, handler := range handlers {
		if handler.filter.Match(HandlerFilter{Type: m.block.Type, FSType: m.block.Fstype}) {
			handler := handler.new()
			handler.Init(m)
			logs.WithF(handler.GetFields()).Trace("New handler")

			m.handlers = append(m.handlers, handler)

			// TODO load configuration for handler
			// TODO if disabled, remove
		}
	}

	// init children managers
	for _, child := range m.block.Children {
		manager := Manager{}
		manager.Init(child, hdm)
		m.children[child.Path] = manager
	}
}

func (m *Manager) HandleEvent(eventType EventType) {
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

func (m *Manager) Start() error {
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

func (m *Manager) Stop(error) {
	close(m.stop)
}
