package rpc

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/password"
)

type HdmServer struct {
	PasswordService *password.Service
}

func (h *HdmServer) SetPassword(password *[]byte, res *struct{}) error {
	return h.PasswordService.FromBytes(password)
}


func (h *HdmServer) Hello(fail *bool, reply *string) error {
	if *fail == true {
		return errs.WithF(data.WithField("test", "test"), "Hello failed!")
	} else {
		*reply = "world"
	}
	return nil
}
