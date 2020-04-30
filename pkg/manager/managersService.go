package manager

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/pilebones/go-udev/netlink"
	"sync"
)

type ManagersService struct {
	PassService *password.Service
	Udev *system.UdevService

	stop              chan struct{}
	managers          map[string]*DiskManager
	managersMutex     sync.RWMutex
}

func (m *ManagersService) Init() {
	m.managers = make(map[string]*DiskManager)
}

func (m *ManagersService) Start() error {
	m.stop = make(chan struct{})
	udevChan := m.Udev.Watch("")
	defer m.Udev.Unwatch(udevChan)

	m.handleEvents(udevChan)

	// cleanup
	m.managersMutex.Lock()
	defer m.managersMutex.Unlock()
	for k, v := range m.managers {
		v.Stop(nil)
		delete(m.managers, k)
	}
	return nil
}

func (m *ManagersService) Stop(e error) {
	close(m.stop)
}

func (m *ManagersService) Get(path string) *DiskManager {
	m.managersMutex.RLock()
	defer m.managersMutex.RUnlock()

	return m.managers[path]
}

///////////////////////

func (m *ManagersService) handleEvents(channel <-chan system.BlockDeviceEvent) {
	for {
		select {
		case event := <-channel:
			if !(event.Action == netlink.ADD && event.Type == "disk")  {
				continue
			}
			m.AddBlockDevice(event)
		case <-m.stop:
			return
		}
	}
}

//

//func (m *ManagersService) handleBlockDeviceEvent(event system.BlockDeviceEvent) {
	//switch event.Action {
	//case "add":
	//	m.AddBlockDevice(event)
	//case "remove":
	//	m.RemoveBlockDevice(event.Path)
	//case "change":
	//	manager := m.Get(event.Path)
	//	if manager == nil {
	//		logs.WithField("path", event.Path).Warn("Receiving change event for unknown block. Creating")
	//		m.AddBlockDevice(event)
	//	} else {
	//		// TODO notify change
	//	}
	//default:
	//	logs.WithField("event", event).Warn("Unknown udev event action")
	//}
//}

//func (m *ManagersService) RemoveBlockDevice(path string) {
//	m.managersMutex.Lock()
//	defer m.managersMutex.Unlock()
//
//	if diskManager, ok := m.managers[path]; ok {
//		diskManager.Stop(nil)
//		delete(m.managers, path)
//	} else {
//		logs.WithField("path", path).Warn("Cannot remove disk, not found")
//	}
//}

func (m *ManagersService) AddBlockDevice(event system.BlockDeviceEvent) {
	m.managersMutex.RLock()
	if _, ok := m.managers[event.Path]; !ok {
		m.managersMutex.RUnlock()

		manager := DiskManager{}

		// TODO handle partitions
		lsblk := hdm.HDM.Servers.GetLocal().Lsblk
		if err := manager.Init(&lsblk, event.Path, m.Udev); err != nil {
			logs.WithE(err).Error("Failed to init blockdevice manager")
			return
		}
		m.register(&manager)
		// TODO add vs start
		manager.HandleEvent(Add)
	} else {
		logs.WithField("path", event.Path).Warn("Cannot register disk, already exists")
	}
}

func (m *ManagersService) register(manager *DiskManager) {
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
