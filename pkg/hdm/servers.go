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

	var inited []Server
	for i := range *s {
		current := (*s)[i]
		if err := current.Init(localHostname); err != nil {
			logs.WithEF(err, data.WithField("server", current)).Warn("Failed to init server")
			continue
		}
		inited = append(inited, current)
	}

	*s = inited

	if len(*s) == 0 {
		*s = append(*s, s.GetLocal())
	}
	return nil
}

func (s *Servers) GetLocal() Server {
	var localServer *Server
	if localServer == nil {
	}

	for _, v := range *s {
		if v.isLocal {
			return v
		}
	}

	localHostname, err := os.Hostname()
	if err != nil {
		logs.WithE(err).Warn("Failed to get local hostname to setup local server")
	}

	logs.WithField("hostname", localHostname).Warn("Local server not found in hdm configuration, creating empty")
	localServer = &Server{}
	localServer.Name = localHostname
	if err := localServer.Init(""); err != nil {
		logs.WithE(err).Error("Failed to init local server properly")
	}
	// TODO this is stupid
	//servers := append(*s, *localServer)
	//*s = servers

	return *localServer
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

		disks, err := srv.Lsblk.ListBlockDevices()
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
