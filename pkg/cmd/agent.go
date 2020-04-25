package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/manager"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/rpc"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
)

func agentCommand() *cobra.Command {
	selector := hdm.DisksSelector{}
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run agent that handle disks lifecycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runningAsRoot(); err != nil {
				return err
			}

			var g run.Group

			//sigterm
			sigterm := utils.SigtermService{}
			sigterm.Init()
			g.Add(sigterm.Start, sigterm.Stop)

			//password
			passService := password.Service{}
			passService.Init()
			//TODO remove
			//pass := []byte("ss")
			//passService.FromBytes(&pass)
			g.Add(passService.Start, passService.Stop)

			//udevService
			udevService := system.UdevService{}
			lsblk := system.Lsblk{}
			if err := lsblk.Init(runner.Local); err != nil {
				return err
			}
			udevService.Init(&lsblk)
			g.Add(udevService.Start, udevService.Stop)

			//manager
			manager := manager.ManagersService{
				PassService: &passService,
				Udev: &udevService,
			}
			manager.Init()
			g.Add(manager.Start, manager.Stop)


			//
			//hdm := rpc.HdmServer{}
			//rpcServer := rpc2.NewServer()
			//if err := rpcServer.Register(&hdm); err != nil {
			//	return errs.WithE(err, "Failed to register hdm rpc server")
			//}
			//// http
			//httpServer := rpc.HttpServer{}
			//httpServer.Init(rpcServer)
			//g.Add(httpServer.Start, httpServer.Stop)
			//
			//// socket
			//socketServer := rpc.SocketServer{}
			//socketServer.Init(rpcServer)
			//g.Add(socketServer.Start, socketServer.Stop)

			rpcServer := rpc.Server{}
			rpcServer.Init(3636, &passService)
			g.Add(rpcServer.Start, rpcServer.Stop)

			// start services
			if err := g.Run(); err != nil {
				return err
			}

			logs.Info("Bye !")

			return nil
		},
	}

	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")

	return cmd
}
