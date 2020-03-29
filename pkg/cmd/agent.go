package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/handlers"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/password"
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
			pass := []byte("ss")
			passService.FromBytes(&pass)
			g.Add(passService.Start, passService.Stop)

			//managers
			managers := handlers.ManagersService{PassService: &passService}
			managers.Init()
			g.Add(managers.Start, managers.Stop)

			//udevService
			udevService := system.UdevService{
				EventChan: managers.GetBlockDeviceEventChan(),
				Filter: selector.Disk,
			}
			lsblk := system.Lsblk{}
			if err := lsblk.Init(runner.Local); err != nil {
				return err
			}
			udevService.Init(&lsblk)
			g.Add(udevService.Start, udevService.Stop)

			/////
			//hdm := rpc.HdmServer{}
			//rpcServer := rpc2.NewServer()
			//if err := rpcServer.Register(&hdm); err != nil {
			//	return errs.WithE(err, "Failed to register hdm rpc server")
			//}
			//
			//// http
			//httpServer := rpc.HttpServer{}
			//httpServer.Init(rpcServer)
			//g.Add(httpServer.Start, httpServer.Stop)
			//
			// socket
			//socketServer := rpc.SocketServer{}
			//socketServer.Init(rpcServer)
			//g.Add(socketServer.Start, socketServer.Stop)

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
