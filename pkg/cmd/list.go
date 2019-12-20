package cmd

import (
	"fmt"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

func listCommand(parent *cobra.Command) {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "list",
		Aliases: []string{"ls", "get"},
		Short: "list disks info",
		Run: errorLoggerWrap(func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "server\tDisk\tLocation\tsize\tLabels"); err != nil { // \tPath,Mount
				logs.WithE(err).Fatal("fail")
			}


			err := hdm.HDM.Servers.RunForDisks(selector, func(srv hdm.Server, disk system.BlockDevice) error {
				locationPath, err := disk.LocationPath()
				if err != nil {
					return err
				}

				location := srv.BayLocation(locationPath)

				var labels []string
				for _, partitions := range disk.Children {
					if partitions.Partlabel != "" {
						labels = append(labels, partitions.Partlabel)
					}
				}

				if _, err := fmt.Fprintln(w,
					srv.Name+"\t"+
					disk.Name+"\t"+
						location+"\t"+
						disk.Size+"\t"+
						strings.Join(labels, ",")+"\t"+
						""); err != nil {
					logs.WithE(err).Fatal("Fail to print")
				}
				return nil
			})
			if err != nil {
				return err
			}
			_ = w.Flush()
			return nil
		}),
	}

	withDiskSelector(&selector, cmd)

	parent.AddCommand(cmd)
}
