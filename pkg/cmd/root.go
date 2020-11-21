package cmd

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/checksum"
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

func RootCmd(version string) *cobra.Command {
	var hdmHome string
	var logLevel string

	config := &config.GlobalConfig{}
	hdm2 := &hdm.Hdm{}

	cmd := &cobra.Command{
		Use:           "hdm",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					return err
				}
				logs.SetLevel(level)
			}

			if err := config.Init(hdmHome); err != nil {
				return errs.WithE(err, "Failed to init hdm")
			}
			return nil
		},
	}

	cmd.AddCommand(
		checksum.RootCmd(config),
		agentCommand(config),
		versionCommand(hdm2, version),
	)

	cmd.PersistentFlags().StringVarP(&hdm2.Home, "home", "H", hdm.DefaultHomeFolder(), "home folder")
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")

	return cmd
}
