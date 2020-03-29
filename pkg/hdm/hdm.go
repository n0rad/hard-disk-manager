package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"io/ioutil"
	"os"
	"time"
)

var HDM Hdm
const pathMnt = "/mnt"
const pathDBDisk = "/db-disk"
const pathConfig = "/config.yaml"

type Hdm struct {
	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}
	DefaultMountPath string
	Servers Servers
	//rpc.SocketServer

	dbDisk DBDisk

	fields        data.Fields
	CheckInterval time.Duration
}

func (hdm *Hdm) Init(home string) error {

	configPath := home + pathConfig

	if hdm.DefaultMountPath == "" {
		hdm.DefaultMountPath = pathMnt
	}

	if _, err := os.Stat(configPath); err == nil {
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errs.WithEF(err, data.WithField("file", configPath), "Failed to read configuration file")
		}

		if err = yaml.Unmarshal(file, hdm); err != nil {
			return errs.WithEF(err, data.WithField("file", configPath), "Invalid configuration format")
		}
	}

	if err := hdm.Servers.Init(); err != nil {
		return errs.WithE(err, "Failed to init servers")
	}

	if err := hdm.dbDisk.Init(home + pathDBDisk); err != nil {
		return errs.WithE(err, "Failed to init db disk")
	}
	return nil
}


func (hdm Hdm) DBDisk() *DBDisk {
	return &hdm.dbDisk
}
