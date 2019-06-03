package hdm

import (
	"encoding/json"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/sfreiberg/simplessh"
)

type Bay struct {
	Path     string
	Location string
}

type Server struct {
	Name     string
	Hostname string
	Username string
	Bays     []Bay

	fields data.Fields
}

// TODO use it and move runner
func (s *Server) Init() error {
	s.fields = data.WithField("server", s.Name)

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

func (s Server) Exec(cmd string) ([]byte, error) {
	logs.WithFields(s.fields).WithField("cmd", cmd).Debug("Running command on server")
	client, err := simplessh.ConnectWithAgent(s.Hostname, s.Username)
	if err != nil {
		return []byte{}, errs.WithEF(err, data.WithField("hostname", s.Hostname).WithField("username", s.Username), "Fail to ssh to server")
	}
	defer client.Close()

	output, err := client.Exec(cmd)
	logs.WithField("output", string(output)).
		WithField("command", cmd).
		Trace("command output")
	if err != nil {
		return []byte{}, errs.WithEF(err, s.fields.WithField("cmd", cmd).WithField("output", string(output)), "Exec command failed")
	}

	return output, nil
}

func (s Server) ScanDisks() (Disks, error) {
	logs.WithField("server", s.Name).Info("Scan disks")
	var disks Disks
	output, err := s.Exec("sudo lsblk -J -O")
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
