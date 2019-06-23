package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"path"
	"strconv"
	"strings"
)

var filesystems = []string{"ext4", "xfs"}

func (b *BlockDevice) Index() (string, error) {
	if b.Mountpoint == "" {
		return "", errs.WithF(b.fields, "Cannot index, disk is not mounted")
	}
	// todo this should be a stream
	output, err := b.server.Exec("sudo find " + b.Mountpoint + " -type f -printf '%A@ %s %P\n'")
	if err != nil {
		return "", errs.WithEF(err, b.fields, "Failed to run du on blockDevice")
	}
	return string(output), nil
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

func (b *BlockDevice) FindHdmConfigs() ([]HdmConfig, error) {
	var hdmConfigs []HdmConfig
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			configs, err := child.FindHdmConfigs()
			if err != nil {
				return hdmConfigs, err
			}
			hdmConfigs = append(hdmConfigs, configs...)
		}
		return hdmConfigs, nil
	}

	if b.Mountpoint == "" {
		return hdmConfigs, errs.WithF(b.fields, "Disk has not mount point")
	}

	configs, err := b.server.Exec("sudo find " + b.Mountpoint + " -type f -not -path '" + b.Mountpoint + pathBackups + "/*' -name " + hdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithEF(err, b.fields, "Failed to find hdm.yaml files")
	}

	lines := strings.Split(string(configs), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		config := HdmConfig{}
		logs.WithF(b.fields.WithField("path", line)).Debug(hdmYamlFilename + " found")
		if err := config.FillFromFile(*b, line); err != nil {
			return hdmConfigs, err
		}
		hdmConfigs = append(hdmConfigs, config)
	}
	return hdmConfigs, nil
}

func (b *BlockDevice) FindNotBackedUp() ([]string, error) {
	if b.Mountpoint == "" {
		return []string{}, errs.WithF(b.fields, "Cannot index, disk is not mounted")
	}

	output, err := b.server.Exec("sudo find " + b.Mountpoint + " -type d ! -name " + hdmYamlFilename + " -printf '%P\n'")
	if err != nil {
		return []string{}, errs.WithEF(err, b.fields, "Failed to find in mountpoint")
	}

	lines := strings.Split(string(output), "\n")
	notBackedupRoots := make(map[string]bool, len(lines))
	for _, line := range lines {
		notBackedupRoots[line] = true
	}

	for _, line := range lines {
		dir := path.Dir(line)
		if _, ok := notBackedupRoots[dir]; ok {
			delete(notBackedupRoots, "line")
		}
	}

	var res []string
	for notBackedup := range notBackedupRoots {
		res = append(res, notBackedup)
	}

	return res, nil
}
