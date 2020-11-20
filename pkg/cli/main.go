package main

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/cmd"
	"go.uber.org/automaxprocs/maxprocs"
	"math/rand"
	"os"
	"syscall"
	"time"
)

var Version = "0"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		logs.WithE(err).Warn("Failed to set maxprocs")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	if err := syscall.Setpriority(syscall.PRIO_PROCESS, syscall.Getpid(), 19); err != nil {
		logs.WithE(err).Warn("Failed to set process as low priority")
	}

	if err := cmd.RootCmd(Version).Execute(); err != nil {
		logs.WithE(err).Fatal("Command failed")
	}
	os.Exit(0)
}
