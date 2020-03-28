package rpc

import (
	"fmt"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

type ServerOld struct {
	Port       int
	Timeout    time.Duration
	SocketPath string

	commands map[string]func(net.Conn) error
	stop     chan struct{}
	listener net.Listener
}

func (s *ServerOld) Init(port int, passService *password.Service) {
	s.Port = 3636
	s.SocketPath = "/tmp/hdm.sock"
	s.Timeout = 10 * time.Second
	s.commands = make(map[string]func(conn net.Conn) error)
	s.commands["password"] = func (conn net.Conn) error {
		if err := passService.FromConnection(conn); err != nil {
			return nil
		}
		logs.Info("Password changed")
		return nil
	}
}

func (s *ServerOld) Start() error {
	s.cleanupSocket()
	s.stop = make(chan struct{}, 1)
	defer close(s.stop)

	listener, err := net.Listen("unix", s.SocketPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", s.SocketPath), "Failed to listen on socket")
	}
	s.listener = listener
	defer s.cleanupSocket()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stop:
				return nil
			default:
				logs.WithE(err).Error("Failed to accept socket connection")
			}
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *ServerOld) Stop(e error) {
	s.stop <- struct{}{}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

//////////////////////////////

func (s *ServerOld) cleanupSocket() {
	_, err := os.Stat(s.SocketPath)
	if os.IsNotExist(err) {
		return
	}

	if err := syscall.Unlink(s.SocketPath); err != nil {
		logs.WithEF(err, data.WithField("path", s.SocketPath)).Warn("Failed to unlink socket")
	}
}

func (s *ServerOld) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			logs.WithE(err).Warn("Socket Connection closed with error")
		}
	}()

	for {
		if err := conn.SetDeadline(time.Now().Add(s.Timeout)); err != nil {
			logs.WithEF(err, data.WithField("timeout", s.Timeout)).Warn("Failed to set deadline on socket connection")
		}

		command, err := readCommand(conn)
		if err != nil {
			if err != io.EOF {
				logs.WithE(err).Error("Failed to read command on socket")
			}
			return
		}

		commandFunc, ok := s.commands[command]
		if !ok {
			_, _ = fmt.Fprintf(conn, "Unknown command %s\n", command)
			return
		}

		commandFunc(conn)
	}
}

func readCommand(conn net.Conn) (string, error) {
	command := ""
	buffer := make([]byte, 1)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		if n == 0 {
			return "", err
		}

		if string(buffer) == " " || string(buffer) == "\n" {
			return command, nil
		}
		command += string(buffer)
	}
}
