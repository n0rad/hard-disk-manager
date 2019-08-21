package socket

import (
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"io"
	"net"
	"syscall"
	"time"
)

type Server struct {
	Port     int
	Timeout  time.Duration
	commands map[string]func(net.Conn) error
}

func (s *Server) Init(port int) {
	s.Port = 3636
	s.Timeout = 10 * time.Second
	s.commands = make(map[string]func(conn net.Conn) error)
	s.commands["password"] = passwordSocketCommand
}

func (s *Server) Start() error {
	_ = syscall.Unlink("/tmp/hdm.sock")
	listener, err := net.Listen("unix", "/tmp/hdm.sock")
	if err != nil {
		return errs.WithE(err, "Failed to listen socket")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.WithE(err).Error("Failed to accept socket connection")
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {

}

//////////////////////////////

func (s *Server) handleConnection(conn net.Conn) {
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

func passwordSocketCommand(conn net.Conn) error {
	buf := memguard.NewBufferFromReaderUntil(conn, '\n')
	hdm.HDM.SetPassword(buf.Seal())
	return nil
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
