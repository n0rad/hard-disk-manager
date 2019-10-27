package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"time"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{},
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

func (h *HandlerConfig) Init(manager *BlockManager) {
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
	if h.manager.BlockDevice.Mountpoint == "" {
		logs.WithF(h.fields).Debug("cannot scan for config, block device is not mounted")
		return nil
	}

	configs, err := hdm.FindConfigs(h.manager.BlockDevice, h.server)
	if err != nil {
		return err
	}

	for _, e := range configs {
		m := BlockManager{
			Type: "path",
			BlockDevice: h.manager.BlockDevice,
			configPath: e.GetConfigPath(),
			config: e,
		}
		if err := m.Init(); err != nil {
			return err
		}
		h.manager.ManagerService.Register(&m)
	}

	return err
}