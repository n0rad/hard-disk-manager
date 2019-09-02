package utils

import (
	"github.com/n0rad/go-erlog/logs"
	"os"
	"os/signal"
	"syscall"
)

type SigtermService struct {
	cancel chan struct{}
}

func (s *SigtermService) Init() {
	s.cancel = make(chan struct{})
}

func (s SigtermService) Start() error {
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		logs.Info("Received SIGTERM, exiting gracefully...")
	case <-s.cancel:
		break
	}
	return nil
}

func (s SigtermService) Stop(e error) {
	close(s.cancel)
}
