package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"os"
)

type DBDisk struct {
	path string
}

func (db *DBDisk) Init(path string) error {

	if err := os.MkdirAll(path, 0755); err != nil {
		return errs.WithEF(err, data.WithField("path", path), "Failed to create disks db path")
	}

	db.path = path
	return nil
}

func (db *DBDisk) Save(disk system.Disk) error {
	if disk.Serial == "" {
		return errs.WithF(data.WithField("path", disk.Path), "Cannot save, disk has no serial")
	}
	diskYaml, err := yaml.Marshal(disk)
	if err != nil {
		return errs.WithE(err, "Failed to marshal disk")
	}

	filePath := db.path + "/" + disk.Serial + ".yaml"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errs.WithEF(err, data.WithField("file", filePath), "Failed to open disk file")
	}
	defer file.Close()

	if _, err := file.Write(diskYaml); err != nil {
		return errs.WithEF(err, data.WithField("path", filePath), "Failed to write disk yaml to file")
	}
	return nil
}

func (db *DBDisk) LoadDisks(/*servers Servers*/) ([]system.Disk, error) {
	var disks []system.Disk
	pathField := data.WithField("path", db.path)

	files, err := ioutil.ReadDir(db.path)
	if err != nil {
		return disks, errs.WithEF(err, pathField, "Failed to read db directory")
	}

	for _, file := range files {
		disk := system.Disk{}
		filePath := db.path + "/" + file.Name()

		bytes, err := ioutil.ReadFile(filePath)
		fileField := data.WithField("file", filePath)
		if err != nil {
			return disks, errs.WithEF(err, fileField, "Failed to read disk db file")
		}

		if err := disk.PopulateFromBytes(bytes); err != nil {
			return disks, errs.WithEF(err, fileField, "Failed to load disk from db file")
		}

		//disk.Init(servers.GetServer(disk.ServerName))

		disks = append(disks, disk)
	}
	return disks, nil
}
