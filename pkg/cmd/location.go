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

func locationCommand(parent *cobra.Command) {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "location",
		Short: "Get disk location",
		Run: errorLoggerWrap(func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "server\tHandlerName\tLocation\tLabel"); err != nil { // \tPath,Mount
				logs.WithE(err).Fatal("fail")
			}

			err := hdm.HDM.Servers.RunForDisks(selector, func(srv hdm.Server, disk system.BlockDevice) error {
				//location, err := disk.LocationPath()
				//if err != nil {
				//	return err
				//}
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
					srv.Name+"\t"+
					disk.Name+"\t"+
						"xxxxx"+"\t"+
						strings.Join(labels, ",")+"\t"+
						path+"\t"+
					//disk.FindDeepestBlockDevice().Mountpoint+"\t"+
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
		}),
	}

	parent.AddCommand(cmd)
}
