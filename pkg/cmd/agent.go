package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/app"
	"github.com/n0rad/hard-disk-manager/pkg/socket"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"github.com/oklog/run"
)

func Agent() error {
	var g run.Group

	// sigterm
	sigterm := utils.SigtermService{}
	sigterm.Init()
	g.Add(sigterm.Start, sigterm.Stop)

	// agent
	agent := app.Agent{}
	g.Add(agent.Start, agent.Stop)

	// socketServer
	socketServer := socket.Server{}
	socketServer.Init(6363)
	g.Add(socketServer.Start, socketServer.Stop)

	// start services
	if err := g.Run(); err != nil {
		logs.WithError(err).Error("A Service failed")
	}

	logs.Info("Bye !")

	return nil
}
