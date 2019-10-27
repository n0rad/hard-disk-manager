package main

import (
	"fmt"
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/cmd"
	"go.uber.org/automaxprocs/maxprocs"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var Version = "0"
var BuildTime = "1970-01-01_00:00:00_UTC"

func sigQuitThreadDump() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			stacktrace := make([]byte, 2<<20)
			length := runtime.Stack(stacktrace, true)
			fmt.Println(string(stacktrace[:length]))

			_ = ioutil.WriteFile("/tmp/"+strconv.Itoa(os.Getpid())+".dump", stacktrace[:length], 0644)
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}

func main() {
	if os.Getuid() != 0 {
		println("hdm must be run as root")
		os.Exit(1)
	}

	if _, err := maxprocs.Set(); err != nil{
		logs.WithE(err).Warn("Failed to set maxprocs")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	sigQuitThreadDump()

	if err := cmd.RootCommand(Version, BuildTime).Execute(); err != nil {
		logs.WithE(err).Fatal("Failed to process args")
	}
	os.Exit(0)
}