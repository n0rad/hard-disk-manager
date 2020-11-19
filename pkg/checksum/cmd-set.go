package checksum

import (
	"github.com/n0rad/hard-disk-manager/pkg/checksum/integrity"
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/spf13/cobra"
)

func setCommand(conf *config.GlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set integrity",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if err := runCmdForPath(conf, arg, func(pathConf config.PathConfig, d integrity.Directory) func(path string) error {
					return d.Set
				}); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.AddCommand(setWatchCommand(conf))
	return cmd
}
