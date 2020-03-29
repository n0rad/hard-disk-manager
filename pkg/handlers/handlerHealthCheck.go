package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"time"
)

func init() {
	//handlers = append(handlers, handler{
	//	HandlerFilter{Type: "disk"},
	//	func() Handler {
	//		return &HandlerHealthCheck{
	//			CommonHandler: CommonHandler{
	//				handlerName: "healthCheck",
	//			},
	//		}
	//	},
	//})
}

type HandlerHealthCheck struct {
	CommonHandler
	CheckInterval time.Duration
}

func (h *HandlerHealthCheck) Init(manager *BlockManager) {
	h.CommonHandler.Init(manager)

	if h.CheckInterval == 0 {
		h.CheckInterval = 6 * time.Hour
	}
}

func (h *HandlerHealthCheck) Start() {
	h.check()

	ticker := time.NewTicker(h.CheckInterval)
	for {
		select {
		case <- ticker.C:
			logs.WithFields(h.fields).Debug("Time to check disk status")
			h.check()
		case <- h.stop:
			ticker.Stop()
			return
		}
	}
}

func (h *HandlerHealthCheck) check() {
	//TODO init elsewhere
	smartctl := system.Smartctl{}
	if err := smartctl.Init(h.manager.BlockDevice.GetExec(), h.manager.BlockDevice); err != nil {
		logs.WithE(err).Error("Failed to create smartctl")
		return
	}
	result, err := smartctl.All()
	if err != nil {
		logs.WithEF(err, h.fields).Error("Failed to get smartctl info")
		return
	}

	if !result.SmartStatus.Passed {
		logs.WithF(h.fields).Error("Smart status is failed")
	}
}