package hdm

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"os"
)

type ServerArray []Server

type Servers struct {
	ServerArray
	local *Server
}

func (s *Servers) Init() error {
	localHostname, err := os.Hostname()
	if err != nil {
		return errs.WithE(err, "Failed to get hostname")
	}

	for i, srv := range s.ServerArray {
		if err := s.ServerArray[i].Init(localHostname); err != nil {
			return err
		}

		if localHostname == srv.Name {
			s.local = &srv
			break
		}
	}

	if s.local == nil {
		logs.WithField("hostname", localHostname).Warn("Local server not found hdm configuration, creating empty")
		s.local = &Server{}
		if err := s.local.Init(""); err != nil {
			return errs.WithE(err, "Failed to init empty server")
		}
	}

	return nil
}

func (s *Servers) GetLocal() Server {
	return *s.local
}

func (s *Servers) GetBlockDeviceByLabel(label string) (system.BlockDevice, error) {
	for _, v := range s.ServerArray {
		// TODO access private
		device, err := v.Lsblk.GetBlockDeviceByLabel(label)
		if err == nil {
			return device, nil
		}
	}
	return system.BlockDevice{}, errs.WithF(data.WithField("label", label), "Block device with label not found")
}
