package system

import (
	"github.com/n0rad/go-erlog/errs"
	"strconv"
	"strings"
)

func (b *BlockDevice) SpaceAvailable() (int, error) {
	output, err := b.exec.ExecShellGetStdout("df " + b.Path + " --output=avail | tail -n +2")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to run 'df' on blockDevice")
	}

	size, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields.WithField("output", string(output)), "Failed to parse 'df' result")
	}
	return size, nil
}

func (b *BlockDevice) InodeUsed() (int, error) {
	output, err := b.exec.ExecShellGetStdout("df " + b.Path + " --output=iused | tail -n +2")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to run 'df' on blockDevice")
	}

	size, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields.WithField("output", string(output)), "Failed to parse 'df' result")
	}
	return size, nil
}