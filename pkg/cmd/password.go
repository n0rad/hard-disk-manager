package cmd

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"github.com/spf13/cobra"
	"io"
	"net"
	"time"
)

func passwordCmd() *cobra.Command {
	var socket string
	var confirm bool
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Send decryption password to local agent using unix socket",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sendPassword(socket, confirm); err != nil {
				logs.WithE(err).Error("Failed to send password")
			}

		},
	}

	cmd.Flags().StringVarP(&socket, "socket", "s", "/tmp/hdm.sock", "Socket")
	cmd.Flags().BoolVarP(&confirm, "confirm", "c", false, "Confirm password")
	return cmd
}

func sendPassword(socketPath string, confirm bool) error {
	passService := password.Service{}
	go passService.Start()
	defer passService.Stop(nil)

	if err := passService.FromStdin(confirm); err != nil {
		return errs.WithE(err, "Failed to ask password")
	}

	// connect
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("socketPath", socketPath), "Failed to connect to socketPath")
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return errs.WithE(err, "Failed to set deadline")
	}

	if err := writeBytes(conn, []byte("password ")); err != nil {
		return errs.WithE(err, "Failed to write command")
	}

	if err := passService.Write(conn); err != nil {
		return errs.WithE(err, "Failed to write key")
	}
	return nil
}

func writeBytes(conn io.Writer, bytes []byte) error {
	var total, written int
	var err error
	for total = 0; total < len(bytes); total += written {
		written, err = conn.Write(bytes[total:])
		if err != nil {
			return err
		}
	}
	return nil
}
