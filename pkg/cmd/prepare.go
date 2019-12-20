package cmd

import (
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/spf13/cobra"
)

func prepareCommand(parent *cobra.Command) {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare new disk with partitions, crypt, mount, ...",
		Run: errorLoggerWrap(func(cmd *cobra.Command, args []string) error {
			return hdm.HDM.Servers.RunForDisks(selector, func(srv hdm.Server, d system.BlockDevice) error {
				//if len(d.Children) != 0 {
				//	return errs.WithF(d.fields, "Cannot prepare disk, some partitions exists")
				//}
				//
				//logs.WithFields(d.fields.WithField("label", label)).Info("Prepare disk")
				//
				//_, err := d.server.Exec("sudo sgdisk -og " + d.Path)
				//if err != nil {
				//	return errs.WithEF(err, d.fields, "Fail to clear partition table")
				//}
				//
				//_, err = d.server.Exec(`sudo sgdisk -n 1:0:0 -t 1:CA7D7CCB-63ED-4C53-861C-1742536059CC -c 1:"` + label + `" ` + d.Path)
				//if err != nil {
				//	return errs.WithEF(err, d.fields, "Fail to create partition")
				//}
				//
				//if err := d.Scan(); err != nil {
				//	return errs.WithEF(err, d.fields, "Fail to rescan disk after luksFormat")
				//}
				//
				//if len(d.Children) != 1 {
				//	return errs.WithF(d.fields, "Number of partitions is not one after prepare")
				//}
				//
				//if _, err = d.server.Exec("echo -n '" + cryptPassword.String() + "' | sudo cryptsetup --verbose --hash=sha512 --cipher=aes-xts-benbi:sha512 --key-size=512 luksFormat " + d.Children[0].Path + " -"); err != nil {
				//	return errs.WithEF(err, d.fields, "Fail to crypt partition")
				//}
				//
				//if err := d.Scan(); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to rescan disk after luksFormat")
				//}
				//
				//if err := d.Children[0].luksOpen(cryptPassword); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to open crypt partition")
				//}
				//
				//if err := d.Scan(); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to rescan disk after luksOpen")
				//}
				//
				//if _, err = d.server.Exec("sudo mkfs.xfs -L " + label + " -f " + d.Children[0].Children[0].Path); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to make filesystem")
				//}
				//
				//if err := d.Scan(); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to rescan disk after luksOpen")
				//}
				//
				//if err := d.Children[0].Children[0].luksClose(); err != nil {
				//	return errs.WithEF(err, d.fields, "Failed to close partition")
				//}
				//
				return nil
			})
		}),
	}

	withRequiredServerDiskAndLabelSelector(&selector, cmd)

	parent.AddCommand(cmd)
}
