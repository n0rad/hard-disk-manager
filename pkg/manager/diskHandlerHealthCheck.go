package manager

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

func init() {
	diskHandlerBuilders["health-check"] = diskHandlerBuilder{
		new: func() DiskHandler {
			return &HandlerHealthCheck{}
		},
	}
}

type HandlerHealthCheck struct {
	CommonDiskHandler
	CheckInterval time.Duration
}

func (h *HandlerHealthCheck) Init(name string, manager *DiskManager) {
	h.CommonBlockHandler.Init(name, &manager.BlockManager)

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
		case <-h.stopChan:
			ticker.Stop()
			return nil
		}
	}
}

func (h *HandlerHealthCheck) Add() error {
	//TODO init elsewhere
	smartctl := system.Smartctl{}
	if err := smartctl.Init(h.manager.block.GetExec(), h.manager.block); err != nil {
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
