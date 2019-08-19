package hdm

import (
	"bufio"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	system "github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
)

type Hdm struct {
	DBPath     string
	Servers    system.Servers
	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}

	fields data.Fields
}

const pathDB = "/db"

func (hdm *Hdm) InitFromFile(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("file", configPath), "Failed to read configuration file")
	}

	if err = yaml.Unmarshal(file, hdm); err != nil {
		return errs.WithEF(err, data.WithField("file", configPath), "Invalid configuration format")
	}

	if hdm.DBPath == "" {
		hdm.DBPath = filepath.Dir(configPath) + pathDB
	}

	if err := hdm.Servers.Init(); err != nil {
		return errs.WithE(err, "Failed to init servers")
	}

	//hdm.disks, err = LoadDisksFromDB(hdm.DBPath, hdm.Servers)
	//if err != nil {
	//	return errs.WithEF(err, data.WithField("path", hdm.DBPath), "Failed to load disks from DB")
	//}
	return nil
}

func (hdm *Hdm) Password() error {
	client, err := rpc.Dial("unix", "/run/hdm.socket")
	if err != nil {
		return errs.WithE(err, "Failed to dial rpc socket")
	}

	in := bufio.NewReader(os.Stdin)
	for {
		line, _, err := in.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		var reply bool
		err = client.Call("Listener.GetLine", line, &reply)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (hdm *Hdm) FindConfigs(b system.BlockDevice) ([]Config, error) {
	var hdmConfigs []Config
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			configs, err := hdm.FindConfigs(child)
			if err != nil {
				return hdmConfigs, err
			}
			hdmConfigs = append(hdmConfigs, configs...)
		}
		return hdmConfigs, nil
	}

	if b.Mountpoint == "" {
		return hdmConfigs, errs.WithF(hdm.fields, "Disk has not mount point")
	}

	configs, err := b.ExecShell("find " + b.Mountpoint + " -type f -not -path '" + b.Mountpoint + pathBackups + "/*' -name " + hdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithEF(err, hdm.fields, "Failed to find hdm.yaml files")
	}

	lines := strings.Split(string(configs), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		config := Config{}
		logs.WithF(hdm.fields.WithField("path", line)).Debug(hdmYamlFilename + " found")
		if err := config.FillFromFile(b, line); err != nil {
			return hdmConfigs, err
		}
		hdmConfigs = append(hdmConfigs, config)
	}
	return hdmConfigs, nil
}
