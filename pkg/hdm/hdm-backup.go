package hdm

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	system "github.com/n0rad/hard-disk-manager/pkg/system"
	"path"
	"strings"
)

func (hdm *Hdm) Backupable(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)

	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		dd := disk.FindDeepestBlockDevice()

		paths, err := hdm.FindNotBackedUp(dd)
		if err != nil {
			return errs.WithEF(err, fields, "Failed to find non backup dirs")
		}
		for _, path := range paths {
			println(path)
		}
		return nil
	})
}


func (hdm *Hdm) Backups() error {
	return nil
}

func (hdm *Hdm) Backup(selector system.DisksSelector) error {
	fields := data.WithField("selector", selector)

	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		configs, err := hdm.FindConfigs(*disk.BlockDevice)
		if err != nil {
			return errs.WithEF(err, fields, "Cannot backup, Failed to load hdm configs files")
		}

		for _, config := range configs {
			if err := config.RunBackups(disks); err != nil {
				return err
			}
		}
		return nil
	})
}

func (hdm *Hdm) FindNotBackedUp(b system.BlockDevice) ([]string, error) {
	if b.Mountpoint == "" {
		return []string{}, errs.WithF(hdm.fields, "Cannot Find Not backed-up, disk is not mounted")
	}

	output, err := b.ExecShell("find " + b.Mountpoint + " -type d -print0 | while read -d $'\\0' dir; do ls -1 \"$dir/"+ hdmYamlFilename +"\"&> /dev/null || echo $dir; done")
	if err != nil {
		return []string{}, errs.WithEF(err, hdm.fields, "Failed to find in mountpoint")
	}

	lines := strings.Split(string(output), "\n")
	notBackedupRoots := make(map[string]bool, len(lines))
	process := make(map[string]bool, len(lines))
	for _, line := range lines {
		//println(line)
		notBackedupRoots[line] = true
		process[line] = true
	}

	for _, line := range lines {
		dir := path.Dir(line)
		//println(dir)
		if _, ok := notBackedupRoots[dir]; ok {
			//println("delete " + line)
			delete(process, line)
		}
	}

	var res []string
	for notBackedup := range process {
		res = append(res, notBackedup)
	}

	return res, nil
}
