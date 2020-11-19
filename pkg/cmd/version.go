package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func versionCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hdm")
			fmt.Println("version : ", version)
			return nil
		},
	}
	return cmd
}
