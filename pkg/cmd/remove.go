package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/spf13/cobra"
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

			if err := runningAsRoot(); err != nil {
				return err
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
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			if err := remove(child); err != nil {
				logs.WithE(err).Warn("Cannot remove device")
			}
		}
	}

	// TODO remove all mount points
	if b.Mountpoint != "" {
		logs.WithFields(b.GetFields()).Info("Umount")
		if err := b.Umount("/mnt/"+b.GetUsableLabel()); err != nil {
			return err
		}
	}

	switch b.Type {
	case "crypt":
		logs.WithFields(b.GetFields()).Info("Luks close")
		if err := b.LuksClose(); err != nil {
			return err
		}
	case "disk":
		logs.WithFields(b.GetFields()).Info("Put in sleep")
		if err := b.PutInSleepNow(); err != nil {
			return err
		}
	}
	return nil
}
