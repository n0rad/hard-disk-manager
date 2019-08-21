package cmd

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/app"
	"github.com/n0rad/hard-disk-manager/pkg/socket"
	"os"
	"os/signal"
	"syscall"
)

func Agent() error {
	// TODO get passwords for disks

	agent := app.Agent{}

	socketServer := socket.Server{}
	socketServer.Init(6363)
	if err := socketServer.Start(); err != nil {
		return errs.WithE(err, "Failed to start socket server")
	}

	if err := agent.Start(); err != nil {
		return errs.WithE(err, "Failed to init agent")
	}

	// Signal handler to quit properly monitor mode
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-signals
	logs.Warn("Shutdown signal received")
	agent.Stop()
	os.Exit(0)

	return nil
}
