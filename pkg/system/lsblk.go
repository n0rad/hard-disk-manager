package system

import (
	"encoding/json"
	"github.com/Masterminds/semver"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
)

var lsblkMinVersion = semver.MustParse("2.33")

type LsblkResult struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
}

type Lsblk struct {
	exec runner.Exec
}

func (l *Lsblk) Init(exec runner.Exec) error {
	if exec == nil {
		return errs.With("Exec cannot be null")
	}
	l.exec = exec

	lsblkVersion, err := l.Version()
	if err != nil {
		return errs.WithE(err, "Failed to get lsblk version to check compatibility")
	}
	if lsblkVersion.LessThan(lsblkMinVersion) {
		return errs.WithF(data.WithField("current", lsblkVersion.String()).WithField("required", lsblkMinVersion.String()), "lsblk version is not compatible with hdm")
	}

	return nil
}

func (l Lsblk) Version() (semver.Version, error) {
	cmd := `lsblk --version | sed "s/.* \([0-9]\+.*\)/\1/"`
	versionString, err := l.exec.ExecShellGetStdout(cmd)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("cmd", cmd), "Failed to get lsblk version")
	}
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("versionString", versionString), "Failed to parse lsblk version")
	}
	return *version, nil
}

func (l Lsblk) GetBlockDevice(path string) (BlockDevice, error) {
	if path == "" {
		return BlockDevice{}, errs.With("Path is required to get blockDevice")
	}

	output, err := l.exec.ExecGetStdout("lsblk", "-J", "-O", path)
	if err != nil {
		return BlockDevice{}, errs.WithEF(err, data.WithField("path", path).WithField("out", output), "Fail to get disk from lsblk")
	}

	lsblk := LsblkResult{}
	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return BlockDevice{}, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	if len(lsblk.Blockdevices) != 1 {
		return BlockDevice{}, errs.WithF(data.WithField("output", output), "Scan disk give more than disk")
	}

	lsblk.Blockdevices[0].Init(l.exec)

	return lsblk.Blockdevices[0], nil
}

func (l Lsblk) ListFlatBlockDevices() ([]BlockDevice, error) {
	lsblk := struct {
		Blockdevices []BlockDevice `json:"blockdevices"`
	}{}

	output, err := l.exec.ExecGetStdout("lsblk", "-J", "-l", "-O", "-e", "2")
	if err != nil {
		return lsblk.Blockdevices, errs.WithE(err, "Fail to get disks from lsblk")
	}

	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return lsblk.Blockdevices, errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	for _, b := range lsblk.Blockdevices {
		b.Init(l.exec)
	}

	return lsblk.Blockdevices, nil
}

func (l Lsblk) GetBlockDeviceByLabel(label string) (BlockDevice, error) {
	devices, err := l.ListFlatBlockDevices()
	if err != nil {
		return BlockDevice{}, err
	}
	for _, device := range devices {
		if device.Label == label {
			device.Init(l.exec)
			return device, nil
		}
	}
	return BlockDevice{}, errs.WithF(data.WithField("label", label), "No block device found with label")
}
