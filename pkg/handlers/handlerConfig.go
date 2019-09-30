package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"time"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "part"},
		func() Handler {
			return &HandlerConfig{
				CommonHandler: CommonHandler{
					handlerName: "handlerConfig",
				},
			}
		},
	})
}

type HandlerConfig struct {
	CommonHandler
	CheckInterval time.Duration
}

func (h *HandlerConfig) Init(manager *BlockDeviceManager) {
	h.CommonHandler.Init(manager)

	if h.CheckInterval == 0 {
		h.CheckInterval = 72 * time.Hour
	}
}

func (h *HandlerConfig) Start() {
	ticker := time.NewTicker(h.CheckInterval)
	if err := h.scan(); err != nil {
		logs.WithEF(err, h.fields).Error("Failed to scan configs")
	}

	for {
		select {
		case <- ticker.C:
			logs.WithFields(h.fields).Debug("Time to scan hdm configs")
			if err := h.scan(); err != nil {
				logs.WithEF(err, h.fields).Error("Failed to scan configs")
			}
		case <- h.stop:
			ticker.Stop()
			return
		}
	}
}

func (h *HandlerConfig) scan() error {
	disk, err := h.server.ScanDisk(h.manager.Path)
	if err != nil {
		return err
	}

	if disk.Mountpoint == "" {
		logs.WithF(h.fields).Debug("cannot scan for config, block device is not mounted")
		return nil
	}

	configs, err := hdm.FindConfigs(disk.Mountpoint, h.server)
	if err != nil {
		return err
	}

	logs.WithF(h.fields).WithField("configs", configs).Warn("Configs")
	return err
}