package main

import (
	"fmt"
	"github.com/n0rad/go-erlog/errs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/gomake"
	"github.com/n0rad/hard-disk-manager/pkg/cmd"
	"github.com/wangjia184/sortedset"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"strings"
)

const buildAssetFolder = "dist/bindata/assets"

func main() {
	gomake.ProjectBuilder().
		WithName("hdm").
		WithStep(&gomake.StepBuild{
			PreBuildHook: func(build gomake.StepBuild) error {
				if err := os.MkdirAll("dist/bindata/assets", 0777); err != nil {
					return err
				}

				// TODO man pages
				//	err := doc.GenMarkdownCustom(cmd.RootCmd("0"), os.Stdout, func(string) string {
				//		return ""
				//	})
				//	if err != nil {
				//		return errs.WithE(err, "Failed to generate doc")
				//	}
				//	return nil

				if err := writeCommitLog(); err != nil {
					return err
				}

				return nil
			},

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

func writeCommitLog() error {
	w, err := os.OpenFile(buildAssetFolder+cmd.PathChangelog, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return errs.WithE(err, "Failed to open commit log")
	}
	defer w.Close()

	repo, err := git.PlainOpen("./")
	if err != nil {
		return err
	}
	cIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return err
	}

	set := sortedset.New()
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.NumParents() > 1 { // Ignore merge commits
			return nil
		}
		set.AddOrUpdate(c.ID().String(), sortedset.SCORE(c.Committer.When.Unix()), c)
		return nil
	})
	if err != nil {
		return err
	}

	currentDate := ""
	s := set.PopMin()
	for s != nil {
		c := s.Value.(*object.Commit)
		msgOneLine := strings.SplitN(c.Message, "\n", 2)[0]
		commitDate := c.Committer.When.Format("## 2006-01-02")
		if commitDate != currentDate {
			fmt.Fprintln(w)
			fmt.Fprintln(w, commitDate)
			currentDate = commitDate
		}
		fmt.Fprint(w, "- ")
		fmt.Fprintln(w, msgOneLine)

		s = set.PopMin()
	}
	return nil
}
