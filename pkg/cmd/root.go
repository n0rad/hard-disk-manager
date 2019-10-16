package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

func RootCommand(Version string, BuildTime string) *cobra.Command {
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
				logs.WithE(err).Fatal("Cannot start, failed to load configuration")
			}
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hdm")
			fmt.Println("version : ", Version)
			fmt.Println("Build Time : ", BuildTime)
		},
	})

	cmd.AddCommand(passwordCmd())

	cmd.AddCommand(&cobra.Command{
		Use:   "agent",
		Short: "Run an agent that self handle disks",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Agent(); err != nil {
				logs.WithE(err).Fatal("Command failed")
			}
		},
	})

	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")
	cmd.PersistentFlags().StringVarP(&hdmHome, "home", "H", homeDotConfigPath()+"/hdm", "configFile")
	return cmd
}

func homeDotConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		logs.WithError(err).Fatal("Failed to find user home folder")
	}
	return home + "/.config"
}
