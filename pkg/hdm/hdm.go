package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"io/ioutil"
	"os"
	"time"
)

var HDM Hdm

const pathMnt = "/mnt"
const pathDBDisk = "/db-disk"
const pathDBBackup = "/db-backup"
const pathConfig = "/config.yaml"

type Hdm struct {
	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}
	DefaultMountPath string
	Servers          Servers
	//rpc.SocketServer

	password *password.Service

	diskDB   DiskDB
	backupDB BackupDB
	// dbFile

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

	if err := hdm.diskDB.Init(home + pathDBDisk); err != nil {
		return errs.WithE(err, "Failed to init db disk")
	}

	if err := hdm.backupDB.Init(home + pathDBBackup, hdm.Servers); err != nil {
		return errs.WithE(err, "Failed to init db backup")
	}

	hdm.password = &password.Service{}
	hdm.password.Init()

	return nil
}

func (hdm Hdm) DiskDB() *DiskDB {
	return &hdm.diskDB
}

func (hdm Hdm) BackupDB() *BackupDB {
	return &hdm.backupDB
}

func (hdm Hdm) PassService() *password.Service {
	return hdm.password
}

func (hdm *Hdm) Start() error {
	return hdm.password.Start()
}

func (hdm *Hdm) Stop(error) {
	hdm.password.Stop(nil)
}
