package system

import (
	"encoding/json"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/tools"
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

func (s Server) ScanDisk(path string) (Disk, error) {
	logs.WithField("server", s.Name).Debug("Scan disks")
	if path == "" {
		return Disk{}, errs.With("Path is reuiqred to scan disk")
	}

	output, err := s.Exec("lsblk", "-J", "-O", path)
	if err != nil {
		return Disk{}, errs.WithE(err, "Fail to get disk from lsblk")
	}

	lsblk := Lsblk{}
	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return Disk{}, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	if len(lsblk.Blockdevices) != 1 {
		return Disk{}, errs.WithF(data.WithField("output", output), "Scan disk give more than disk")
	}

	lsblk.Blockdevices[0].Init(&s)

	return lsblk.Blockdevices[0], nil
}

func (s Server) ListFlatBlockDevices() ([]BlockDevice, error) {
	logs.WithField("server", s.Name).Debug("List block devices")
	lsblk := struct {
		Blockdevices []BlockDevice `json:"blockdevices"`
	}{}

	output, err := s.Exec("lsblk", "-J", "-l", "-O")
	if err != nil {
		return lsblk.Blockdevices, errs.WithE(err, "Fail to get disks from lsblk")
	}

	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return lsblk.Blockdevices, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	return lsblk.Blockdevices, nil
}
