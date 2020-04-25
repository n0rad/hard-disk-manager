package system

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/pilebones/go-udev/netlink"
	"strings"
	"sync"
)

// Types: part, lvm, crypt, dmraid, mpath, path, dm, loop, md, linear, raid0, raid1, raid4, raid5, raid10, multipath, disk, tape, printer, processor, worm, rom, scanner, mo-disk, changer, comm, raid, enclosure, rbc, osd, and no-lun

type BlockDeviceEvent struct {
	Action netlink.KObjAction
	Path   string
	Type   string
	FSType string
}

type UdevService struct {
	stop                       chan struct{}
	lsblk                      *Lsblk
	watchersWithPathFilter     map[chan BlockDeviceEvent]string
	watchersWithPathFilterLock sync.RWMutex
}

func (s *UdevService) Init(lsblk *Lsblk) {
	s.lsblk = lsblk
}

func (s *UdevService) Start() error {
	s.stop = make(chan struct{})
	s.watchersWithPathFilter = make(map[chan BlockDeviceEvent]string)

	udevConn := new(netlink.UEventConn)
	if err := udevConn.Connect(netlink.UdevEvent); err != nil {
		return errs.WithE(err, "Unable to connect to Netlink Kobject UEvent socket")
	}
	defer udevConn.Close()

	//if err := s.addCurrentBlockDevices(); err != nil {
	//	s.Stop(err)
	//	return errs.WithE(err, "Cannot add current block devices after watching events")
	//}

	// TODO you can lose events between addCurrent and watch but watch is blocking
	s.watchUdevBlockEvents(udevConn)


	s.watchersWithPathFilterLock.Lock()
	defer s.watchersWithPathFilterLock.Unlock()
	for channel := range s.watchersWithPathFilter {
		close(channel)
	}
	return nil
}

func (s *UdevService) Stop(e error) {
	close(s.stop)
}

// TODO the return chan should be unidirectional but it requires to change input param type for unwatch
func (s *UdevService) Watch(filter string) chan BlockDeviceEvent {
	s.watchersWithPathFilterLock.Lock()
	defer s.watchersWithPathFilterLock.Unlock()

	c := make(chan BlockDeviceEvent)
	s.watchersWithPathFilter[c] = filter
	return c
}

func (s *UdevService) Unwatch(c chan BlockDeviceEvent) {
	s.watchersWithPathFilterLock.Lock()
	defer s.watchersWithPathFilterLock.Unlock()

	close(c)
	delete(s.watchersWithPathFilter, c)
}

///////////////////////////////

//func (s *UdevService) addCurrentBlockDevices() error {
//	blockDevices, err := s.lsblk.ListFlatBlockDevices()
//	if err != nil {
//		return errs.WithE(err, "Failed to list current block devices")
//	}
//	for _, v := range blockDevices {
//		if !strings.Contains(v.Path, s.Filter) {
//			continue
//		}
//		s.EventChan <- BlockDeviceEvent{
//			Action: netlink.ADD,
//			Type:   v.Type,
//			Path:   "/dev/" + v.Kname,
//			FSType: v.Fstype,
//		}
//	}
//	return nil
//}

func (s *UdevService) watchUdevBlockEvents(udevConn *netlink.UEventConn) {
	matcher := netlink.RuleDefinitions{
		Rules: []netlink.RuleDefinition{
			{
				Env: map[string]string{
					"SUBSYSTEM": "block",
				},
			},
		},
	}

	queue := make(chan netlink.UEvent)
	defer close(queue)
	errors := make(chan error)
	defer close(errors)
	quitMonitor := udevConn.Monitor(queue, errors, &matcher)
	for {
		select {
		case uevent := <-queue:
			logs.WithField("uevent", uevent).Trace("Received udev event")

			// make udev compatible with lsblk
			if uevent.Env["DEVTYPE"] == "partition" {
				uevent.Env["DEVTYPE"] = "part"
			}

			// TODO mapper device /dev/mdX to /dev/mapper
			//path := uevent.Env["DEVNAME"]
			//if device, err := k.server.GetBlockDevice(path); err != nil {
			//	logs.WithE(err).Warn("Failed to get blockdevice from kernel event")
			//} else {
			//	// replace kernel path with lsblk path (/dev/dmX -> /dev/mapper/XX)
			//	path = device.Path
			//}

			event := BlockDeviceEvent{
				Action: uevent.Action,
				Path:   uevent.Env["DEVNAME"],
				Type:   uevent.Env["DEVTYPE"],
				FSType: uevent.Env["ID_FS_TYPE"],
			}

			s.watchersWithPathFilterLock.RLock()
			for channel, filter := range s.watchersWithPathFilter {
				if !strings.Contains(uevent.Env["DEVNAME"], filter) {
					continue
				}
				channel <- event
			}
			s.watchersWithPathFilterLock.RUnlock()
		case err := <-errors:
			logs.WithE(err).Warn("Received error for udev")
		case <-s.stop:
			close(quitMonitor)
			return
		}
	}
}
