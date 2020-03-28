package rpc

import (
	"bytes"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"io"
	"net"
	"net/http"
	rpc2 "net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type HttpServer struct {
	Address string
	Timeout time.Duration

	rpcServer *rpc2.Server
	stop      chan struct{}
	listener  net.Listener
}

func (s *HttpServer) Init(rpcServer *rpc2.Server) {
	if s.Address == "" {
		s.Address = ":8686"
	}
	if s.Timeout == 0 {
		s.Timeout = 10 * time.Second
	}
	s.rpcServer = rpcServer
}

func (s HttpServer) Start() error {
	s.stop = make(chan struct{}, 1)
	defer close(s.stop)


	logs.WithField("address", s.Address).Info("Listen http server")
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", s.Address), "Failed to listen on socket")
	}
	s.listener = listener

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call(s.rpcServer)
		_, _ = io.Copy(w, res)
	})
	//
	if err := http.Serve(s.listener, mux); err != nil {
		logs.WithE(err).Warn("Failed on http server")
	}
	return nil
}

func (s *HttpServer) Stop(e error) {
	s.stop <- struct{}{}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

/////////////////////////////////////////

type rpcRequest struct {
	r    io.Reader     // holds the JSON formated RPC request
	rw   io.ReadWriter // holds the JSON formated RPC response
	done chan bool     // signals then end of the RPC request
}

// NewRPCRequest returns a new rpcRequest.
func NewRPCRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

// Read implements the io.ReadWriteCloser Read method.
func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *rpcRequest) Call(server *rpc2.Server) io.Reader {
	go server.ServeCodec(jsonrpc.NewServerCodec(r))
	//go server.ServeConn(r)
	//jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}
