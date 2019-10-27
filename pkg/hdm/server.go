package hdm

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type Bay struct {
	Path     string
	Location string
}

type Server struct {
	Name          string
	Hostname      string
	LocalHostname string
	Username      string
	Bays          []Bay

	fields data.Fields
	Lsblk  system.Lsblk // TODO public ?
}

func (s *Server) Init(localHostname string) error {
	s.fields = data.WithField("server", s.Name)

	var exec runner.Exec
	if localHostname == s.Name || localHostname == "" || localHostname == "localhost" {
		exec = &runner.LocalExec{}
	} else {
		exec = &runner.SshExec{
			Username: s.Username,
			Hostname: s.LocalHostname, // TODO
		}
	}
	s.Lsblk = system.Lsblk{}
	if err := s.Lsblk.Init(exec); err != nil {
		return errs.WithEF(err, s.fields, "Failed to init lsblk for server")
	}

	return nil
}

func (s *Server) BayLocation(path string) string {
	for _, bay := range s.Bays {
		if bay.Path == path {
			return bay.Location
		}
	}
	return ""
}
