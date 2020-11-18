package trash

import (
	"fmt"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"os"
	"strings"
	"text/tabwriter"
)

//func List() error {
//	disks, err := system.LoadDisksFromDB(hdm.DBPath, hdm.Servers)
//	if err != nil {
//		logs.Fatal("Failed to load disks from DB")
//	}
//
//	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
//	if _, err := fmt.Fprintln(w, "Label\tRota\tTran\tSize\tserver\tdays\tuncorrect"); err != nil {
//		logs.WithE(err).Fatal("fail")
//	}
//	for _, disk := range disks {
//		if _, err := fmt.Fprintln(w, disk.Label+"\t"+
//			strconv.FormatBool(disk.Rota)+"\t"+
//			disk.Tran+"\t"+
//			disk.Size+"\t"+
//			disk.ServerName+"\t"+
//			strconv.Itoa(disk.SmartResult.PowerOnTime.Hours/24)+"\t"); err != nil {
//			logs.WithE(err).Fatal("fail")
//		}
//	}
//	_ = w.Flush()
//	return nil
//}

func Index(selector system.DisksSelector) error {
	return hdm.HDM.Servers.RunForDisks(selector, func(disks system.Disks, disk system.Disk) error {
		//res, err := findDeepestBlockDevice(disk.BlockDeviceOLD).Index()
		//if err != nil {
		//	return err
		//}
		//print(res)
		//return err
		return nil
	})
}
