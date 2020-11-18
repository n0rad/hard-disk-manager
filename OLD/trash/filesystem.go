package trash

import (
	"github.com/n0rad/go-erlog/errs"
	"strconv"
	"strings"
)

var filesystems = []string{"ext4", "xfs"}

func (b *BlockDeviceOLD) Index() (string, error) {
	if b.Mountpoint == "" {
		return "", errs.WithF(b.fields, "Cannot index, disk is not mounted")
	}
	// todo this should be a stream
	output, err := b.server.Exec("sudo find " + b.Mountpoint + " -type f -printf '%A@ %s %P\n'")
	if err != nil {
		return "", errs.WithEF(err, b.fields, "Failed to find files in filesystem")
	}
	return string(output), nil
}

func (b *BlockDeviceOLD) SpaceAvailable() (int, error) {
	output, err := b.server.Exec("df " + b.Path + " --output=avail | tail -n +2")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to run 'df' on blockDevice")
	}

	size, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields.WithField("output", string(output)), "Failed to parse 'df' result")
	}
	return size, nil
}
