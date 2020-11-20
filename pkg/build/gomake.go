package main

import (
	"github.com/n0rad/go-erlog/errs"
	_ "github.com/n0rad/go-erlog/register"
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
			PostReleaseHook: func(step gomake.StepRelease) error {
				if err := gomake.Exec("docker", "build", ".",
					"-f", "Dockerfile.release",
					"-t", "n0rad/hdm:"+step.Version,
					"-t", "n0rad/hdm:latest"); err != nil {
					return errs.WithE(err, "Failed to build docker image")
				}

				if err := gomake.Exec("docker", "push", "n0rad/hdm:"+step.Version); err != nil {
					return errs.WithE(err, "Failed to push docker image")
				}

				if err := gomake.Exec("docker", "push", "n0rad/hdm:latest"); err != nil {
					return errs.WithE(err, "Failed to push latest docker image")
				}

				return nil
			},
		}).
		MustBuild().MustExecute()
}
