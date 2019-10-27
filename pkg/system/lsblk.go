package system

import (
	"github.com/Masterminds/semver"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
)

type Lsblk struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
}

func LsblkVersion() (semver.Version, error) {
	cmd := `lsblk --version | sed "s/.* \(2.*\)/\1/"`
	versionString, err := runner.Local.ExecShellGetStdout(cmd)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("cmd", cmd), "Failed to get lsblk version")
	}
	version, err := semver.NewVersion(versionString)
	if err != nil {
		return semver.Version{}, errs.WithEF(err, data.WithField("versionString", versionString), "Failed to parse lsblk version")
	}
	return *version, nil
}