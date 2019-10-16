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

func (s *Server) ExecShell(command string) (string, error) {
	return s.Exec("sh", "-c", command)
}


func (s Server) GetBlockDevice(path string) (BlockDevice, error) {
	if path == "" {
		return BlockDevice{}, errs.With("Path is required to get blockDevice")
	}

	output, err := s.Exec("lsblk", "-J", "-O", path)
	if err != nil {
		return BlockDevice{}, errs.WithEF(err, data.WithField("path", path), "Fail to get disk from lsblk")
	}

	lsblk := Lsblk{}
	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return BlockDevice{}, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	if len(lsblk.Blockdevices) != 1 {
		return BlockDevice{}, errs.WithF(data.WithField("output", output), "Scan disk give more than disk")
	}

	lsblk.Blockdevices[0].Init(&s)

	return lsblk.Blockdevices[0], nil
}

func (s Server) ListFlatBlockDevices() ([]BlockDevice, error) {
	logs.WithField("server", s.Name).Debug("List block devices")
	lsblk := struct {
		Blockdevices []BlockDevice `json:"blockdevices"`
	}{}

	output, err := s.Exec("lsblk", "-J", "-l", "-O", "-e", "2")
	if err != nil {
		return lsblk.Blockdevices, errs.WithE(err, "Fail to get disks from lsblk")
	}

	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return lsblk.Blockdevices, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	return lsblk.Blockdevices, nil
}

func (s Server) GetBlockDeviceByLabel(label string) (BlockDevice, error) {
	devices, err := s.ListFlatBlockDevices()
	if err != nil {
		return BlockDevice{}, err
	}
	for _, device := range devices {
		if device.Label == label {
			return device, nil
		}
	}
	return BlockDevice{}, errs.WithF(data.WithField("label", label),"No block device found with label")
}
