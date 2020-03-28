package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func versionCommand(version string, buildTime string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hdm")
			fmt.Println("version : ", version)
			fmt.Println("build time : ", buildTime)
			return nil
		},
	}
	return cmd
}
