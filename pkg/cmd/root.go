package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/checksum"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var configFile string
	//var config = &Config{}

	var logLevel string
	cmd := &cobra.Command{
		Use:           "hdm",
		SilenceErrors: true,
		SilenceUsage:  true,
		//PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//	return config.Load(configFile)
		//},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					logs.WithField("value", logLevel).Fatal("Unknown log level")
				}
				logs.SetLevel(level)
			}
		},
	}

	cmd.AddCommand(
		checksum.RootCmd(),
	)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", `./hdm.yaml`, "configuration file")
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")

	cmd.MarkFlagRequired("config")

	return cmd
}
