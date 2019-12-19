package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type DisksSelector struct {
	Server string
	Label  string
	Disk   string
}

func (d DisksSelector) MatchDisk(srv Server, disk system.BlockDevice) bool {
	if d.Server == srv.Name || d.Server == "" {
		if d.Disk == "" && d.Label == "" {
			return true
		}
		if d.Disk != "" && d.Disk == disk.Name {
			return true
		}
		if d.Label != "" && d.Label == disk.Partlabel {
			return true
		}
		for _, child := range disk.Children {
			res := d.MatchPartition(srv, disk, child)
			if res {
				return true
			}
		}
		return false
	}
	return false
}

func (d DisksSelector) MatchPartition(srv Server, disk system.BlockDevice, device system.BlockDevice) bool {
	if d.Server != srv.Name {
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