package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/handlers"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/n0rad/hard-disk-manager/pkg/socket"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"github.com/oklog/run"
)

func Agent() error {
	var g run.Group

	// sigterm
	sigterm := utils.SigtermService{}
	sigterm.Init()
	g.Add(sigterm.Start, sigterm.Stop)

	// password
	passService := password.Service{}
	passService.Init()
	g.Add(passService.Start, passService.Stop)

	// managers
	managers := handlers.ManagersService{PassService: &passService}
	managers.Init()
	g.Add(managers.Start, managers.Stop)

	// udevService
	udevService := system.UdevService{
		EventChan: managers.GetBlockDeviceEventChan(),
	}
	udevService.Init()
	g.Add(udevService.Start, udevService.Stop)

	// socketServer
	socketServer := socket.Server{}
	socketServer.Init(6363, &passService)
	g.Add(socketServer.Start, socketServer.Stop)

	// start services
	if err := g.Run(); err != nil {
		return err
	}

	logs.Info("Bye !")

	return nil
}
