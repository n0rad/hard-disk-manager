package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"strconv"
	"strings"
)

func (b *BlockDevice) Index() error {
	//sudo find /mnt/2000 -type f -printf "%A@ %s %P\n"
	//find ./ -type f -printf "%A@ %s %f\n"
	return nil
}

func (b *BlockDevice) SpaceAvailable() (int, error) {
	output, err := b.server.Exec("df " + b.Path + " --output=avail | tail -n +2")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to run du on blockDevice")
	}

	size, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields.WithField("output", string(output)), "Failed to parse 'df' result")
	}
	return size, nil
}