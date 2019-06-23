package hdm

import (
	"bufio"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-drive-manager/pkg/utils"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
)

type Hdm struct {
	DBPath  string
	Servers Servers
	LuksFormat []struct {
		Hash string
		Cipher string
		keySize string
	}
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

	if err := hdm.Servers.init(); err != nil {
		return errs.WithE(err, "Failed to init servers")
	}

	//hdm.disks, err = LoadDisksFromDB(hdm.DBPath, hdm.Servers)
	//if err != nil {
	//	return errs.WithEF(err, data.WithField("path", hdm.DBPath), "Failed to load disks from DB")
	//}
	return nil
}

func (hdm *Hdm) List() error {
	disks, err := LoadDisksFromDB(hdm.DBPath, hdm.Servers)
	if err != nil {
		logs.Fatal("Failed to load disks from DB")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "Label\tRota\tTran\tSize\tserver\tdays\tuncorrect"); err != nil {
		logs.WithE(err).Fatal("fail")
	}
	for _, disk := range disks {
		if _, err := fmt.Fprintln(w, disk.Label+"\t"+
			strconv.FormatBool(disk.Rota)+"\t"+
			disk.Tran+"\t"+
			disk.Size+"\t"+
			disk.ServerName+"\t"+
			strconv.Itoa(disk.SmartResult.PowerOnTime.Hours/24)+"\t"); err != nil {
			logs.WithE(err).Fatal("fail")
		}
	}
	_ = w.Flush()
	return nil
}

func (hdm *Hdm) Index(selector DisksSelector) error {
	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		//res, err := findDeepestBlockDevice(disk.BlockDevice).Index()
		//if err != nil {
		//	return err
		//}
		//print(res)
		//return err
		return nil
	})
}

func (hdm *Hdm) Add(selector DisksSelector) error {
	password, err := utils.AskPasswordWithConfirmation(false)
	if err != nil {
		return errs.WithE(err, "Failed to get password")
	}

	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		return disk.Add(password)
	})
}

func (hdm *Hdm) Remove(selector DisksSelector) error {
	fields := data.WithField("selector", selector)

	disk, err := hdm.Servers.GetDisk(selector)
	if err != nil {
		return err
	}
	if disk == nil {
		return errs.WithF(fields, "Disk not found")
	}

	return disk.Remove()
}

func (hdm *Hdm) Location(selector DisksSelector) error {
	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		location, err := disk.Location()
		if err != nil {
			return err
		}
		println(location)
		return nil
	})
}

func (hdm *Hdm) Prepare(selector DisksSelector) error {
	fields := data.WithField("selector", selector)
	label := selector.Label
	selector.Label = ""

	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		if disk.HasChildren() {
			return errs.WithF(fields, "Cannot prepare, disk has partitions")
		}

		password, err := utils.AskPasswordWithConfirmation(true)
		if err != nil {
			return errs.WithE(err, "Failed to get password")
		}

		return disk.Prepare(label, password)
	})
}

func (hdm *Hdm) Backupable(selector DisksSelector) error {
	fields := data.WithField("selector", selector)

	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		dd := disks.findDeepestBlockDeviceByLabel(selector.Label) // TODO that sux hard
		if dd == nil {
			return errs.WithF(fields, "disk not found")
		}

		paths, err := dd.FindNotBackedUp()
		if err != nil {
			return errs.WithEF(err, fields, "Failed to find non backup dirs")
		}
		for _, path := range paths {
			println(path)
		}
		return nil
	})
}

func (hdm *Hdm) Backup(selector DisksSelector) error {
	fields := data.WithField("selector", selector)

	return hdm.Servers.RunForDisks(selector, func(disks Disks, disk Disk) error {
		configs, err := disk.FindHdmConfigs()
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

func (hdm *Hdm) Agent() error {
	// get passwords for disks
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
