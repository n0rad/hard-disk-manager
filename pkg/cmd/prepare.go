package cmd

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/spf13/cobra"
	"os"
)

func prepareCommand() *cobra.Command {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare new disk with partitions, crypt, mount, ...",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := hdm.HDM.Servers.GetLocal().Lsblk.GetBlockDevice(selector.Disk)
			if err != nil {
				return err
			}

			if os.Getuid() != 0 {
				return errs.With("Being root required")
			}

			if d.HasChildren() {
				return errs.WithF(data.WithField("block", d), "Cannot prepare disk, some children exists")
			}

			logs.WithFields(data.WithField("block", d)).Info("Prepare disk")

			logs.WithFields(data.WithField("block", d)).Info("Clear partition table disk")
			if err := d.ClearPartitionTable(); err != nil {
				return err
			}

			if err := d.Reload(); err != nil {
				return errs.WithEF(err, data.WithField("block", d), "Fail to reload disk info")
			}

			logs.WithFields(data.WithField("block", d)).Info("Create partition")
			if err := d.CreateSinglePartition(selector.Label); err != nil {
				return err
			}

			if err := d.Reload(); err != nil {
				return errs.WithEF(err, data.WithField("block", d), "Fail to reload disk info")
			}

			if len(d.Children) != 1 {
				return errs.WithF(data.WithField("block", d), "Number of partitions is not one after creation")
			}

			passService := password.Service{}
			passService.Init()
			go passService.Start()
			defer passService.Stop(nil)

			if err := passService.FromStdin(true); err != nil {
				return errs.WithE(err, "Failed to ask password")
			}

			pass, err := passService.Get()
			if err != nil {
				return errs.WithE(err, "Failed to get password from lock storage")
			}

			logs.WithFields(data.WithField("block", d.Children[0])).Info("Encrypt partition")
			if err := d.Children[0].LuksFormat(pass); err != nil {
				return err
			}

			if err := d.Reload(); err != nil {
				return errs.WithEF(err, data.WithField("block", d), "Fail to reload disk info")
			}

			logs.WithFields(data.WithField("block", d.Children[0])).Info("Open Encrypted partition")
			if err := d.Children[0].LuksOpen(pass); err != nil {
				return err
			}

			if err := d.Reload(); err != nil {
				return errs.WithEF(err, data.WithField("block", d), "Fail to reload disk info")
			}

			if len(d.Children[0].Children) != 1 {
				return errs.WithF(data.WithField("block", d.Children[0]), "Number of children block device after open is not one")
			}

			logs.WithFields(data.WithField("block", d.Children[0].Children[0])).Info("Create filesystem")
			if err := d.Children[0].Children[0].Format("xfs", selector.Label); err != nil {
				return err
			}

			logs.WithFields(data.WithField("block", d.Children[0].Children[0])).Info("Close encrypted partition")
			if err := d.Children[0].Children[0].LuksClose(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")
	cmd.Flags().StringVarP(&selector.Label, "label", "l", "", "Label")

	_ = cmd.MarkFlagRequired("disk")
	_ = cmd.MarkFlagRequired("label")

	return cmd
}
