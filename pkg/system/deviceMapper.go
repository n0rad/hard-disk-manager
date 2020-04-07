package system

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"os"
)

//for file in $(ls -la /dev/mapper/* | grep "\->" | grep -oP "\-> .+" | grep -oP " .+"); do echo "MAPPER:"$(F=$(echo $file | grep -oP "[a-z0-9-]+");echo $F":"$(ls "/sys/block/${F}/slaves/");)":"$(df -h "/dev/mapper/${file}" | sed 1d); done;
// dmsetup table
// dmsetup remove /dev/mapper/yopla

// find who use the mount point : lsof /mnt/yopla
// find if still mounted : mount | grep
// remove device mapper : dmsetup remove


//func findFromBlockDevice(name string) {
//
//}


type DeviceMapper struct {
	Exec runner.Exec
}

func (d DeviceMapper) BlockFromName(name string) (string, error) {
	// TODO this is not compatible with remote exec
	s, err := os.Readlink("/dev/mapper/"+name)
	if err != nil {
		return "", errs.WithE(err,"Failed to get link from devicemapper")
	}
	return s, nil
}

func (d DeviceMapper) Remove(name string) error {
	std, err := d.Exec.ExecGetStd("dmsetup", "remove", "/dev/mapper/"+name)
	if err != nil {
		return errs.WithEF(err, data.WithField("std", std).WithField("name", name), "Failed to remove device mapper")
	}
	return nil
}