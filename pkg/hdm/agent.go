package hdm

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"os"
	"os/signal"
	"syscall"
)

func (hdm *Hdm) Agent() error {
	// TODO get passwords for disks

	agent := Agent{}

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
