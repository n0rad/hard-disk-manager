package manager

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"sync"
)

type ManagersService struct {
	PassService *password.Service
	Udev *system.UdevService

	blockDeviceEvents chan system.BlockDeviceEvent
	stop              chan struct{}
	managers          map[string]*DiskManager
	managersMutex     sync.RWMutex
}

func (m *ManagersService) Init() {
	m.blockDeviceEvents = make(chan system.BlockDeviceEvent)
}

func (m *ManagersService) GetBlockDeviceEventChan() chan<- system.BlockDeviceEvent {
	return m.blockDeviceEvents
}

func (m *ManagersService) Start() error {
	m.managers = make(map[string]*DiskManager)
	m.stop = make(chan struct{})

	m.handleEvents()

	// cleanup
	m.managersMutex.Lock()
	defer m.managersMutex.Unlock()
	for _, v := range m.managers {
		v.Stop(nil)
	}
	return nil
}

func (m *ManagersService) Stop(e error) {
	close(m.stop)
}

///////////////////////

func (m *ManagersService) handleEvents() {
	for {
		select {
		case e := <-m.blockDeviceEvents:
			logs.WithField("event", e).Debug("Received block event")
			m.handleBlockDeviceEvent(e)
		case <-m.stop:
			return
		}
	}
}

func (m *ManagersService) Register(manager *DiskManager) {
	m.managersMutex.Lock()
	defer m.managersMutex.Unlock()

	if manager == nil {
		logs.Warn("Trying to register nil manager")
		return
	}

	// TODO handle return
	go manager.Start()

	m.managers[manager.block.Path] = manager
}

//

func (m *ManagersService) Get(path string) *DiskManager {
	m.managersMutex.RLock()
	defer m.managersMutex.RUnlock()

	return m.managers[path]
}

func (m *ManagersService) handleBlockDeviceEvent(event system.BlockDeviceEvent) {
	switch event.Action {
	case "add":
		m.AddBlockDevice(event)
	case "remove":
		m.RemoveBlockDevice(event.Path)
	case "change":
		manager := m.Get(event.Path)
		if manager == nil {
			logs.WithField("path", event.Path).Warn("Receiving change event for unknown block. Creating")
			m.AddBlockDevice(event)
		} else {
			// TODO notify change
		}
	default:
		logs.WithField("event", event).Warn("Unknown udev event action")
	}
}

func (m *ManagersService) RemoveBlockDevice(path string) {
	m.managersMutex.Lock()
	defer m.managersMutex.Unlock()

	if diskManager, ok := m.managers[path]; ok {
		diskManager.Stop(nil)
		delete(m.managers, path)
	} else {
		logs.WithField("path", path).Warn("Cannot remove disk, not found")
	}
}

func (m *ManagersService) AddBlockDevice(event system.BlockDeviceEvent) {
	m.managersMutex.RLock()
	if _, ok := m.managers[event.Path]; !ok {
		m.managersMutex.RUnlock()

		manager := DiskManager{}

		// TODO handle partitions
		if err := manager.Init(hdm.HDM.Servers.GetLocal().Lsblk, event.Path, m.Udev); err != nil {
			logs.WithE(err).Warn("Failed to init blockdevice manager")
		}
		m.Register(&manager)
	} else {
		logs.WithField("path", event.Path).Warn("Cannot register disk, already exists")
	}
}
