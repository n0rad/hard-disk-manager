package system

import (
	"bufio"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"os"
	"strings"
)

type Mount struct {
	Device     string
	Path       string
	Filesystem string
	Flags      string
}

const procSelfMounts = "/proc/self/mounts"

func Mounts() ([]Mount, error) {
	file, err := os.Open(procSelfMounts)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mounts := []Mount(nil)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 5)
		if len(parts) != 5 {
			return nil, errs.WithF(data.WithField("line", line), procSelfMounts+" line has not 5 parts")
		}
		mounts = append(mounts, Mount{parts[0], parts[1], parts[2], parts[3]})
	}

	if err := scanner.Err(); err != nil {
		return nil, errs.WithE(err, procSelfMounts+" read error")
	}

	return mounts, nil
}

func MountFromBlockDevice(blockDevicePath string) (*Mount, error) {
	mounts, err := Mounts()
	if err != nil {
		return nil, errs.WithE(err, "Failed to get mounts")
	}

	for _, mount := range mounts {
		if mount.Device == blockDevicePath {
			return &mount, nil
		}
	}

	return nil, nil
}