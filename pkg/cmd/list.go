package cmd

import (
	"fmt"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

func listCommand() *cobra.Command {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "get"},
		Short:   "list disks info",
		RunE: func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "Server\tSeen\tDisk\tLocation\tSize\tAvailable\tLabels\tHealth"); err != nil { // \tPath,Mount
				return errs.WithE(err, "Fail to prepare writer header")
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
						"12h"+"\t"+
						disk.Name+"\t"+
						location+"\t"+
						disk.Size+"\t"+
						"12G"+"\t"+
						strings.Join(labels, ",")+"\t"+
						"ok"+"\t"+
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
		},
	}

	withDiskSelector(&selector, cmd)

	return cmd
}
