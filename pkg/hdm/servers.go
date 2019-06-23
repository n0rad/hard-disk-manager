package hdm

import (
	"github.com/n0rad/go-erlog/errs"
)

type DisksSelector struct {
	Server string
	Label  string
	Disk   string
}

func (d DisksSelector) MatchDisk(disk Disk) bool {
	if d.Disk != "" && d.Disk == disk.Name {
		return true
	}
	for _, child := range disk.Children {
		res := d.MatchPartition(disk, child)
		if res {
			return true
		}
	}
	return false
}


func (d DisksSelector) MatchPartition(disk Disk, device BlockDevice) bool {
	if d.Server != disk.ServerName {
		return false
	}
	if d.Label != "" {
		if d.Label == device.Partlabel {
			return true
		}
		return false
	} else if d.Disk != "" {
		if d.Disk == disk.Name {
			return true
		}
		return false
	}
	return true
}

func (d DisksSelector) String() string {
	return "server=" + d.Server + ",disk=" + d.Disk + ",label=" + d.Label
}

func (d DisksSelector) IsValid() error {
	if d.Disk == "" && d.Label == "" {
		return errs.With("disk or label flag are mandatory")
	}
	if d.Disk != "" && d.Label != "" {
		return errs.With("disk and label cannot be set at the same time")
	}
	if d.Disk != "" && d.Server == "" {
		return errs.With("server is mandatory if disk is set")
	}
	return nil
}

type Servers []Server

func (s *Servers) init() error {
	for _, server := range *s {
		if err := server.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (s Servers) GetServer(name string) *Server {
	for _, srv := range s {
		if srv.Name == name {
			return &srv
		}
	}
	return nil
}

func (s Servers) GetDisk(selector DisksSelector) (*Disk, error) {
	for _, srv := range s {
		if srv.Name != selector.Server {
			continue
		}
		disks, err := srv.ScanDisks()
		if err != nil {
			return nil, errs.WithE(err, "Failed to scan disks")
		}
		return disks.findDiskBySelector(selector), nil
	}
	return nil, nil
}

func (s Servers) ScanDisks() ([]Disk, error) {
	var allDisks []Disk
	for _, srv := range s {
		disks, err := srv.ScanDisks()
		if err != nil {
			return allDisks, errs.WithE(err, "Failed to list disks")
		}
		allDisks = append(allDisks, disks...)
	}
	return allDisks, nil
}

func (s Servers) RunForDisks(selector DisksSelector, toRun func(disks Disks, disk Disk) error) error {
	for _, srv := range s {
		if selector.Server != "" && selector.Server != srv.Name {
			continue
		}

		disks, err := srv.ScanDisks()
		if err != nil {
			return err
		}

		for _, disk := range disks {
			if !selector.MatchDisk(disk) {
				continue
			}
			if err := toRun(disks, disk); err != nil {
				return errs.WithE(err, "Command failed")
			}
		}
	}
	return nil
}

///sys/block/sda/device/scsi_device/0\:0\:0\:0/device/

//