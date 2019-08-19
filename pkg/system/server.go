package system

import (
	"encoding/json"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/tools"
	"strings"
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
	runner tools.Runner
}

// TODO use it and move runner
func (s *Server) Init() error {
	s.fields = data.WithField("server", s.Name)
	s.runner = &tools.LocalRunner{}
	//s.runner = &tools.SshRunner{
	//	Hostname: s.Hostname,
	//	Username: s.Username,
	//}
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

func (s Server) Exec(head string, args ...string) (string, error) {
	stdout, _, err := s.runner.ExecGetOutputError(head, args...)
	return stdout, err
}

func (s Server) ScanDisks() (Disks, error) {
	logs.WithField("server", s.Name).Info("Scan disks")
	var disks Disks
	output, err := s.Exec("lsblk", "-J", "-O")
	if err != nil {
		return disks, errs.WithE(err, "Fail to get disks from lsblk")
	}

	lsblk := Lsblk{}
	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return disks, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	for i := range lsblk.Blockdevices {
		lsblk.Blockdevices[i].Init(&s)

		if lsblk.Blockdevices[i].Name == "fd0" {
			logs.WithFields(s.fields.WithField("device", lsblk.Blockdevices[i].Name)).Debug("Skipping device")
			continue
		}

		//if err := lsblk.Blockdevices[i].FillFromSmartctl(); err != nil {
		//	return disks, errs.WithE(err, "Failed to add smartctl info disk")
		//}

		lsblk.Blockdevices[i].ServerName = s.Name

		disks = append(disks, lsblk.Blockdevices[i])
	}

	return disks, nil
}

func (s Server) ListDisks() ([]string, error) {
	logs.WithField("server", s.Name).Debug("List disks")
	output, err := s.Exec("lsblk", "-n", "-d", "-o", "path")
	if err != nil {
		return []string{}, errs.WithE(err, "Fail to get disks from lsblk")
	}
	return strings.Split(string(output), "\n"), nil
}
