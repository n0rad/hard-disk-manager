package main

import (
	"fmt"
	"github.com/n0rad/go-erlog/logs"
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/cmd"
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
	rand.Seed(time.Now().UTC().UnixNano())
	sigQuitThreadDump()

	if err := cmd.RootCommand(Version, BuildTime).Execute(); err != nil {
		logs.WithE(err).Fatal("Failed to process args")
	}
	os.Exit(0)
}
