package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func versionCommand(root *cobra.Command, version string, buildTime string) {
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		Run: errorLoggerWrap(func(cmd *cobra.Command, args []string) error {
			fmt.Println("hdm")
			fmt.Println("version : ", version)
			fmt.Println("build time : ", buildTime)
			return nil
		}),
	})
}
