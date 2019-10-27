package hdm

import (
	"github.com/Masterminds/semver"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"os"
	"time"
)

var HDM Hdm
var lsblkMinVersion = semver.MustParse("2.33")
var smartCtlMinVersion = semver.MustParse("7.0")

type Hdm struct {
	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}

	DefaultMountPath string

	dbDisk DBDisk

	fields        data.Fields
	CheckInterval time.Duration
}

const pathMnt = "/mnt"
const pathDBDisk = "/db-disk"
const pathConfig = "/config.yaml"

func (hdm Hdm) DBDisk() *DBDisk {
	return &hdm.dbDisk
}

func (hdm *Hdm) CheckVersions() error {
	lsblkVersion, err := system.LsblkVersion()
	if err != nil {
		return errs.WithE(err, "Failed to get lsblk version to check compatibility")
	}
	if lsblkVersion.LessThan(lsblkMinVersion) {
		return errs.WithF(data.WithField("current", lsblkVersion.String()).WithField("required", lsblkMinVersion.String()), "lsblk version is not compatible with hdm")
	}

	// smartctl
	smartctlVersion, err := system.SmartctlVersion()
	if err != nil {
		return errs.WithE(err, "Failed to get smartctl version to check compatibility")
	}
	if smartctlVersion.LessThan(smartCtlMinVersion) {
		logs.WithF(data.WithField("current", smartctlVersion.String()).WithField("required", smartCtlMinVersion.String())).Error("smartctl version is not compatible with hdm")
	}

	return nil
}

func (hdm *Hdm) Init(home string) error {
	if err := hdm.CheckVersions(); err != nil {
		return err
	}

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

	if err := hdm.dbDisk.Init(home + pathDBDisk); err != nil {
		return errs.WithE(err, "Failed to init db disk")
	}
	return nil
}
