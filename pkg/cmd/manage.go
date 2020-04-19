package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/handler"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
)

func manageCommand() *cobra.Command {
	var diskName string
	cmd := &cobra.Command{
		Use:   "manage",
		Short: "Manage a disk",
	}

	cmd.PersistentFlags().StringVarP(&diskName, "p", "d", "", "disk")
	_ = cmd.MarkFlagRequired("disk")

	cmd.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Manage a disk to add",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runningAsRoot(); err != nil {
				return err
			}

			manage, err := startManage(diskName)
			if err != nil {
				return err
			}

			manage.HandleEvent(handler.Add)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "Manage a disk to remove",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runningAsRoot(); err != nil {
				return err
			}

			manage, err := startManage(diskName)
			if err != nil {
				return err
			}

			manage.HandleEvent(handler.Remove)
			return nil
		},
	})

	return cmd
}

func startManage(diskName string) (*handler.DiskManager, error) {
	var g run.Group

	//sigterm
	sigterm := utils.SigtermService{}
	sigterm.Init()
	g.Add(sigterm.Start, sigterm.Stop)

	// hdm
	g.Add(hdm.HDM.Start, hdm.HDM.Stop)

	// lsblk
	lsblk := system.Lsblk{}
	lsblk.Init(runner.Local)

	// disk manager
	m := handler.DiskManager{}
	if err := m.Init(lsblk, diskName); err != nil {
		return nil, err
	}
	g.Add(m.Start, m.Stop)

	// start services
	go func(g *run.Group) {
		if err := g.Run(); err != nil {
			logs.WithE(err).Error("error")
			//return err
		}
	}(&g)

	return &m, nil
}