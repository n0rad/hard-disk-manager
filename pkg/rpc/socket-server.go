package rpc

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"syscall"
	"time"
)

type SocketServer struct {
	Timeout    time.Duration
	SocketPath string

	rpcServer *rpc.Server
	stop      chan struct{}
	listener  net.Listener
}

func (s *SocketServer) Init(rpcServer *rpc.Server) {
	if s.SocketPath == "" {
		s.SocketPath = "/tmp/hdm.sock"
	}
	if s.Timeout == 0 {
		s.Timeout = 10 * time.Second
	}
	s.rpcServer = rpcServer
}

func (s SocketServer) Start() error {
	s.cleanupSocket()
	s.stop = make(chan struct{}, 1)
	defer close(s.stop)


	logs.WithField("address", s.SocketPath).Info("Listen socket server")
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

		if err := conn.SetDeadline(time.Now().Add(s.Timeout)); err != nil {
			logs.WithEF(err, data.WithField("timeout", s.Timeout)).Warn("Failed to set deadline on socket connection")
		}

		go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}

func (s *SocketServer) Stop(e error) {
	s.stop <- struct{}{}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

/////////////////////////////////////////

func (s *SocketServer) cleanupSocket() {
	_, err := os.Stat(s.SocketPath)
	if os.IsNotExist(err) {
		return
	}

	if err := syscall.Unlink(s.SocketPath); err != nil {
		logs.WithEF(err, data.WithField("path", s.SocketPath)).Warn("Failed to unlink socket")
	}
}
