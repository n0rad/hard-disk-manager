package hdm

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"os"
)

type Servers []Server

func (s *Servers) Init() error {
	localHostname, err := os.Hostname()
	if err != nil {
		return errs.WithE(err, "Failed to get hostname")
	}

	var localServer *Server
	for i := range *s {
		current := &(*s)[i]
		if err := current.Init(localHostname); err != nil {
			return err
		}

		if localHostname == current.Name {
			localServer = current
		}
	}

	if localServer == nil {
		logs.WithField("hostname", localHostname).Warn("Local server not found in hdm configuration, creating empty")

		localServer = &Server{}
		localServer.Name = localHostname
		if err := localServer.Init(""); err != nil {
			return errs.WithE(err, "Failed to init empty server")
		}
		servers := append(*s, *localServer)
		*s = servers
	}

	return nil
}

func (s Servers) GetLocal() Server {
	for _, v := range s {
		if v.isLocal {
			return v
		}
	}
	logs.Error("Illegal state: Get local server found no server")
	return Server{}
}

func (s Servers) GetBlockDeviceByLabel(label string) (system.BlockDevice, error) {
	for _, v := range s {
		// TODO access private
		device, err := v.Lsblk.GetBlockDeviceByLabel(label)
		if err == nil {
			logs.WithE(err).Warn("erf")
			return device, nil
		}
	}
	return system.BlockDevice{}, errs.WithF(data.WithField("label", label), "Block device with label not found")
}

func (s Servers) RunForDisks(selector DisksSelector, toRun func(srv Server, disk system.BlockDevice) error) error {
	for _, srv := range s {
		if selector.Server != "" && selector.Server != srv.Name {
			continue
		}

		disks, err := srv.Lsblk.ListFlatBlockDevices()
		if err != nil {
			return err
		}

		for _, disk := range disks {
			if !selector.MatchDisk(srv, disk) {
				continue
			}
			if err := toRun(srv, disk); err != nil {
				logs.WithE(err).Error("Run for disk Command failed")
				//return errs.WithEF(err, disk.fields,)
			}
		}
	}
	return nil
}
