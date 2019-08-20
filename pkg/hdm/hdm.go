package hdm

import (
	"bufio"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strings"
)

var HDM Hdm

type Hdm struct {
	Servers    system.Servers

	LuksFormat []struct {
		Hash    string
		Cipher  string
		keySize string
	}

	dbDisk DBDisk

	fields data.Fields
}

const pathDBDisk = "/db-disk"
const pathConfig = "/config.yaml"

func (hdm Hdm) DBDisk() *DBDisk {
	return &hdm.dbDisk
}

func (hdm *Hdm) Init(home string) error {
	configPath := home+pathConfig
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

	if err := hdm.dbDisk.Init(home+pathDBDisk); err != nil {
		return errs.WithE(err, "Failed to init db disk")
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

	configs, err := b.ExecShell("find " + b.Mountpoint + " -type f -not -path '" + b.Mountpoint + pathBackups + "/*' -name " + HdmYamlFilename)
	if err != nil {
		return hdmConfigs, errs.WithEF(err, hdm.fields, "Failed to find hdm.yaml files")
	}

	lines := strings.Split(string(configs), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		config := Config{}
		logs.WithF(hdm.fields.WithField("path", line)).Debug(HdmYamlFilename + " found")
		if err := config.FillFromFile(b, line); err != nil {
			return hdmConfigs, err
		}
		hdmConfigs = append(hdmConfigs, config)
	}
	return hdmConfigs, nil
}
