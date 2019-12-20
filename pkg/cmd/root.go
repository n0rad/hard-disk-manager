package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

func RootCommand(version string, buildTime string) *cobra.Command {
	var logLevel string
	var hdmHome string

	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					logs.WithField("value", logLevel).Fatal("Unknown log level")
				}
				logs.SetLevel(level)
			}

			if err := hdm.HDM.Init(hdmHome); err != nil {
				logs.WithE(err).Fatal("Failed to init hdm")
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")
	cmd.PersistentFlags().StringVarP(&hdmHome, "home", "H", homeDotConfigPath()+"/hdm", "configFile")

	versionCommand(cmd, version, buildTime)
	passwordCommand(cmd)
	agentCommand(cmd)
	listCommand(cmd)
	prepareCommand(cmd)
	return cmd
}

func homeDotConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		logs.WithError(err).Fatal("Failed to find user home folder")
	}
	return home + "/.config"
}
