package handler

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

func init() {
	diskHandlerBuilders["info"] = diskHandlerBuilder{
		new: func() DiskHandler {
			return &HandlerInfo{}
		},
	}
}

type HandlerInfo struct {
	CommonDiskHandler
}

///////////////////////////////////
func (h *HandlerInfo) Change() error {
	return h.Add()
}

func (h *HandlerInfo) Add() error {
	if h.manager.block.Serial == "" {
		logs.WithF(h.GetFields()).Trace("Disk has no serial, not saving")
		return nil
	}

	if err := hdm.HDM.DiskDB().Save(h.manager.block); err != nil {
		return errs.WithE(err, "Failed to save disk to db")
	}
	return nil
}
