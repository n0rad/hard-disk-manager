package main

import "github.com/n0rad/gomake"

func main() {
	gomake.ProjectBuilder().
		WithStep(&gomake.StepBuild{
			BinaryName: "hdm",
		}).
		MustBuild().MustExecute()
}
