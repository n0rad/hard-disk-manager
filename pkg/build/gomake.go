//+build build

package main

import (
	"github.com/n0rad/gomake"
)

func main() {
	gomake.ProjectBuilder().
		WithName("hdm").
		WithStep(&gomake.StepBuild{
			Programs: []gomake.Program{
				{
					BinaryName: "hdm",
					Package:    "github.com/n0rad/hard-disk-manager/pkg/cli/hdm",
				},
			},
		}).
		WithStep(&gomake.StepRelease{
			OsArchRelease: []string{"linux-amd64"},
		}).
		MustBuild().MustExecute()
}
