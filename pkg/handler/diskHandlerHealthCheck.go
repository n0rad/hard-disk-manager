package handler

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

func init() {
	DiskHandlerBuilders["health-check"] = diskHandlerBuilder{
		//block.HandlerFilter{Type: "disk"},
		New: func() DiskHandler {
			return &HandlerHealthCheck{
				CommonDiskHandler: CommonDiskHandler{
					CommonBlockHandler: CommonBlockHandler{
						CommonHandler: CommonHandler{
							HandlerName: "health-check",
						},
					},
				},
			}
		},
	}
}

type HandlerHealthCheck struct {
	CommonDiskHandler
	CheckInterval time.Duration
}

func (h *HandlerHealthCheck) Init(manager *DiskManager) {
	h.CommonBlockHandler.Init(&manager.BlockManager)

	if h.CheckInterval == 0 {
		h.CheckInterval = 6 * time.Hour
	}
}

func (h *HandlerHealthCheck) Start() error {
	ticker := time.NewTicker(h.CheckInterval)
	for {
		select {
		case <-ticker.C:
			logs.WithFields(h.GetFields()).Debug("Time to check disk status")
			if err := h.Add(); err != nil {
				logs.WithE(err).Error("Health check failed")
			}
		case <-h.StopChan:
			ticker.Stop()
			return nil
		}
	}
}

func (h *HandlerHealthCheck) Add() error {
	//TODO init elsewhere
	smartctl := system.Smartctl{}
	if err := smartctl.Init(h.manager.Block.GetExec(), h.manager.Block); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to create smartctl")
	}
	result, err := smartctl.All()
	if err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to get smartctl info")
	}

	if !result.SmartStatus.Passed {
		logs.WithF(h.GetFields()).Error("Smart status is failed")
	}

	return nil
}
