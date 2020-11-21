package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/checksum"
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/spf13/cobra"
)

func RootCmd(version string) *cobra.Command {
	var hdmHome string
	var logLevel string

	config := &config.GlobalConfig{}

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
		versionCommand(version),
	)

	homeDotConfig, err := homeDotConfigPath()
	if err != nil {
		homeDotConfig = "/tmp"
	}

	cmd.PersistentFlags().StringVarP(&hdmHome, "home", "H", homeDotConfig+"/hdm", "configFile")
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")

	return cmd
}

func homeDotConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", errs.WithE(err, "Failed to find user home folder")
	}
	return home + "/.config", nil
}
