package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"time"
)

var HDM Hdm

type Hdm struct {
	Servers system.Servers

	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}

	dbDisk   DBDisk

	fields data.Fields
	CheckInterval time.Duration
}

const pathDBDisk = "/db-disk"
const pathConfig = "/config.yaml"

func (hdm Hdm) DBDisk() *DBDisk {
	return &hdm.dbDisk
}

func (hdm *Hdm) Init(home string) error {
	configPath := home + pathConfig
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("file", configPath), "Failed to read configuration file")
	}

	if err = yaml.Unmarshal(file, hdm); err != nil {
		return errs.WithEF(err, data.WithField("file", configPath), "Invalid configuration format")
	}

	if err := hdm.Servers.Init(); err != nil {
		return errs.WithE(err, "Failed to init servers")
	}

	if err := hdm.dbDisk.Init(home + pathDBDisk); err != nil {
		return errs.WithE(err, "Failed to init db disk")
	}
	return nil
}

