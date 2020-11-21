package hdm

import (
	"github.com/juju/fslock"
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/dist"
	"io/ioutil"
	"os"
)

const pathAssets = "/assets"
const pathLock = "/lock"
const pathVersion = "/version"

//var HDM = Hdm{}

type Hdm struct {
	Home    string
	version string

	assetsExtracted bool
}

func (hdm *Hdm) GetAssetsFolder() (string, error) {
	if !hdm.assetsExtracted {
		if err := hdm.RestoreAssets(); err != nil {
			return "", errs.WithE(err, "Failed to restore asset folder")
		}
		hdm.assetsExtracted = true
	}
	return hdm.Home + pathAssets, nil
}

func (hdm *Hdm) RestoreAssets() error {
	if err := os.MkdirAll(hdm.Home, 0755); err != nil {
		return errs.WithE(err, "Failed to mkdir Home directory")
	}
	lock := fslock.New(hdm.Home + pathLock)
	err := lock.Lock()
	if err != nil {
		return errs.WithE(err, "Failed to get asset extract lock")
	}
	defer lock.Unlock()

	bytes, err := ioutil.ReadFile(hdm.Home + pathVersion)
	if err != nil {
		logs.WithE(err).Warn("Failed to read Home version. May be first run")
	}
	if string(bytes) != hdm.version || err != nil {
		logs.WithField("homeVersion", string(bytes)).WithField("currentVersion", hdm.version).Info("Hdm version changed, extract of assets required")

		//
		assetsPath := hdm.Home + pathAssets
		if err := os.RemoveAll(assetsPath); err != nil {
			return errs.WithEF(err, data.WithField("path", assetsPath), "Failed to cleanup old assets")
		}

		if err := dist.RestoreAssets(hdm.Home, "assets"); err != nil {
			return errs.WithEF(err, data.WithField("path", assetsPath), "Failed to restore assets")
		}

		if err := ioutil.WriteFile(hdm.Home+pathVersion, []byte(hdm.version), 0644); err != nil {
			logs.WithE(err).Error("Failed to write current Hdm version to Home")
		}
	}

	return nil
}

func DefaultHomeFolder() string {
	home, err := homedir.Dir()
	if err != nil {
		logs.WithE(err).Warn("Failed to find Home directory")
		home = "/tmp/hdm-Home"
	}
	return home + "/.config/hdm"
}
