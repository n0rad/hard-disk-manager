package cmd

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const PathChangelog = "/CHANGELOG.md"

func versionChangelogCommand(hdm *hdm.Hdm) *cobra.Command {
	return &cobra.Command{
		Use:                   "changelog [oldVersion]",
		Short:                 "Display bbc changelog",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			assets, err := hdm.GetAssetsFolder()
			if err != nil {
				return err
			}
			file, err := os.Open(assets + PathChangelog)
			if err != nil {
				return errs.WithE(err, "Failed to open changelog file")
			}
			defer file.Close()
			if _, err := io.Copy(os.Stdout, file); err != nil {
				return errs.WithE(err, "Failed to print changelog file to stdout")
			}
			return nil
		},
	}
}
