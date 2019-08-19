package hdm

import (
	"fmt"
	"github.com/n0rad/go-erlog/logs"
	system "github.com/n0rad/hard-disk-manager/pkg/system"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func (hdm *Hdm) List() error {
	disks, err := system.LoadDisksFromDB(hdm.DBPath, hdm.Servers)
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

func (hdm *Hdm) Index(selector system.DisksSelector) error {
	return hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		//res, err := findDeepestBlockDevice(disk.BlockDevice).Index()
		//if err != nil {
		//	return err
		//}
		//print(res)
		//return err
		return nil
	})
}


func (hdm *Hdm) Location(selector system.DisksSelector) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "Name\tLocation\tLabel\tPath,Mount"); err != nil {
		logs.WithE(err).Fatal("fail")
	}

	err := hdm.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		location, err := disk.Location()
		if err != nil {
			return err
		}
		path, err := disk.LocationPath()
		if err != nil {
			return err
		}

		var labels []string
		for _, partitions := range disk.Children {
			if partitions.Partlabel != "" {
				labels = append(labels, partitions.Partlabel)
			}
		}


		if _, err := fmt.Fprintln(w,
			disk.Name+"\t"+
				location+"\t"+
				strings.Join(labels, ",")+"\t"+
				path+"\t" +
				disk.FindDeepestBlockDevice().Mountpoint+"\t"+
				""); err != nil {
			logs.WithE(err).Fatal("Fail tp print")
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = w.Flush()
	return nil
}
