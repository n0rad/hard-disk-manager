package main

import (
	"github.com/n0rad/gomake"
)

func main() {
	gomake.ProjectBuilder().
		WithName("hdm").
		WithStep(&gomake.StepBuild{
			// TODO generate readme
			//PreBuildHook: func(build gomake.StepBuild) error {
			//	err := doc.GenMarkdownCustom(cmd.RootCmd("0"), os.Stdout, func(string) string {
			//		return ""
			//	})
			//	if err != nil {
			//		return errs.WithE(err, "Failed to generate doc")
			//	}
			//	return nil
			//},

			Programs: []gomake.Program{
				{
					BinaryName: "hdm",
					Package:    "github.com/n0rad/hard-disk-manager/pkg/cli",
				},
			},
		}).
		WithStep(&gomake.StepRelease{
			OsArchRelease: []string{"linux-amd64"},
		}).
		MustBuild().MustExecute()
}
