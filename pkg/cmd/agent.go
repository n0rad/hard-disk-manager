package cmd

import (
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func agentCommand(conf *config.GlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run as an agent handling disks and files lifecycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			return nil
		},
	}
	return cmd
}
