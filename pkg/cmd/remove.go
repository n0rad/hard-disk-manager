package cmd

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/spf13/cobra"
	"os"
)

func removeCommand() *cobra.Command {
	selector := hdm.DisksSelector{}
	selector.Server = "n02"
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove or cleanup removed disk",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := selector.IsValid(); err != nil {
				return err
			}

			if os.Getuid() != 0 {
				return errs.With("Being root required")
			}

			return hdm.HDM.Servers.RunForDisks(selector, func(srv hdm.Server, disk system.BlockDevice) error {
				return remove(disk)
			})
		},
	}

	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")
	cmd.Flags().StringVarP(&selector.Label, "label", "l", "", "Label")

	return cmd
}

func remove(b system.BlockDevice) error {
	logs.WithFields(b.GetFields()).Info("Disk remove")
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			if err := remove(child); err != nil {
				logs.WithE(err).Warn("Cannot remove device")
			}
		}
	}

	// TODO remove all mount points
	if b.Mountpoint != "" {
		if err := b.Umount("/mnt/"+b.GetUsableLabel()); err != nil {
			return err
		}
	}

	switch b.Type {
	case "crypt":
		if err := b.LuksClose(); err != nil {
			return err
		}
	case "disk":
		if err := b.PutInSleepNow(); err != nil {
			return err
		}
	}
	return nil
}
