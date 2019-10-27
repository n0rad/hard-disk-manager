package system

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/pilebones/go-udev/netlink"
)


// Types: part, lvm, crypt, dmraid, mpath, path, dm, loop, md, linear, raid0, raid1, raid4, raid5, raid10, multipath, disk, tape, printer, processor, worm, rom, scanner, mo-disk, changer, comm, raid, enclosure, rbc, osd, and no-lun

type BlockDeviceEvent struct {
	Action netlink.KObjAction
	Path   string
	Type   string
	FSType string
}

type UdevService struct {
	EventChan chan <-BlockDeviceEvent

	stop   chan struct{}
	server Server
}

func (k *UdevService) Init() error {
	if err := k.server.Init(); err != nil {
		return errs.WithE(err, "Failed to init empty server")
	}

	return nil
}

func (k *UdevService) Start() error {
	k.stop = make(chan struct{})

	udevConn := new(netlink.UEventConn)
	if err := udevConn.Connect(netlink.UdevEvent); err != nil {
		return errs.WithE(err, "Unable to connect to Netlink Kobject UEvent socket")
	}
	defer udevConn.Close()

	if err := k.addCurrentBlockDevices(); err != nil {
		k.Stop(err)
		return errs.WithE(err, "Cannot add current block devices after watching events")
	}

	// TODO you can lose events between addCurrent and watch but watch is blocking
	k.watchUdevBlockEvents(udevConn)

	logs.Info("Stop Agent")

	return nil
}

func (k *UdevService) Stop(e error) {
	close(k.stop)
}

///////////////////////////////

func (k *UdevService) addCurrentBlockDevices() error {
	blockDevices, err := k.server.ListFlatBlockDevices()
	if err != nil {
		return errs.WithE(err, "Failed to list current block devices")
	}
	for _, v := range blockDevices {
		k.EventChan <- BlockDeviceEvent{
			Action: netlink.ADD,
			Type:   v.Type,
			Path:   "/dev/"+v.Kname,
			FSType: v.Fstype,
		}
	}
	return nil
}

func (k *UdevService) watchUdevBlockEvents(udevConn *netlink.UEventConn) {
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

			if uevent.Env["DEVTYPE"] == "partition" {
				uevent.Env["DEVTYPE"] = "part"
			}

			//path := uevent.Env["DEVNAME"]
			//if device, err := k.server.GetBlockDevice(path); err != nil {
			//	logs.WithE(err).Warn("Failed to get blockdevice from kernel event")
			//} else {
			//	// replace kernel path with lsblk path (/dev/dmX -> /dev/mapper/XX)
			//	path = device.Path
			//}

			k.EventChan <- BlockDeviceEvent{
				Action: uevent.Action,
				Path:   uevent.Env["DEVNAME"],
				Type:   uevent.Env["DEVTYPE"],
				FSType: uevent.Env["ID_FS_TYPE"],
			}

		case err := <-errors:
			logs.WithE(err).Warn("Received error for udev watcher")
		case <-k.stop:
			close(quitMonitor)
			return
		}
	}
}
