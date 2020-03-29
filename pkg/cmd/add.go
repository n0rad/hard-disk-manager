package cmd

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

func addCommand() *cobra.Command {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add disk(s) for usage (mdadm, luksOpen, mount, ...)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getuid() != 0 {
				return errs.With("Being root required")
			}

			passService := password.Service{}
			passService.Init()
			go passService.Start()
			defer passService.Stop(nil)

			return hdm.HDM.Servers.RunForDisks(selector, func(srv hdm.Server, disk system.BlockDevice) error {
				if _, err := addAndGiveNewDevices(disk, &passService); err != nil {
					return err
				}
				return nil
			})
		},
	}

	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")
	cmd.Flags().StringVarP(&selector.Label, "label", "l", "", "Label")

	return cmd
}

var filesystems = []string{"ext4", "xfs"}

func addAndGiveNewDevices(d system.BlockDevice, passService *password.Service) (bool, error) {
	logs.WithFields(d.GetFields()).Debug("Add device")

	if len(d.Children) > 0 {
		newDevices := false
		for _, child := range d.Children {
			newRecursive, err := addAndGiveNewDevices(child, passService)
			if err != nil {
				logs.WithEF(err, d.GetFields()).Warn("Cannot add device")
			}
			if newRecursive == true {
				newDevices = newRecursive
			}
		}
		return newDevices, nil
	}

	newDevices := false
	if d.Fstype == "crypto_LUKS" {
		if !passService.IsSet() {
			if err := passService.FromStdin(false); err != nil {
				return false, errs.WithE(err, "Failed to ask password")
			}
		}

		pass, err := passService.Get()
		if err != nil {
			return false, errs.WithE(err, "Failed to get password from lock storage")
		}

		if err := d.LuksOpen(pass); err != nil {
			return false, err
		}
		newDevices = true
	} else if utils.SliceContains(filesystems, d.Fstype) {
		mountPath := "/mnt/" + d.GetUsableLabel()

		if err := d.Mount(mountPath); err != nil {
			d.Umount(mountPath)
			return false, err
		}
	} else {
		return false, errs.WithF(d.GetFields().WithField("fstype", d.Fstype), "Unknown fstype")
	}
	return newDevices, nil
}
