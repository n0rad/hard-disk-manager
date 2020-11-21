package cmd

import (
	"fmt"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

func versionCommand(hdm *hdm.Hdm, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hdm")
			fmt.Println("version : ", version)
			return nil
		},
	}

	cmd.AddCommand(versionChangelogCommand(hdm))

	return cmd
}
