package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"time"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "disk"},
		func() Handler {
			return &HandlerDb{
				CommonHandler: CommonHandler{
					handlerName: "healthCheck",
				},
			}
		},
	})
}

// store disk info in db
type HandlerHealthCheck struct {
	CommonHandler
	CheckInterval time.Duration
}

func (h *HandlerHealthCheck) Init(manager *BlockDeviceManager) {
	h.CommonHandler.Init(manager)

	if h.CheckInterval == 0 {
		h.CheckInterval = 10 * time.Second
	}
}

func (h *HandlerHealthCheck) Start() {
	ticker := time.NewTicker(h.CheckInterval)
	for {
		select {
		case <- ticker.C:
			logs.WithFields(h.fields).Debug("Time to check disk status")

			//err := h.disk.FillFromSmartctl()
			//
			//
			//if !h.disk.SmartResult.SmartStatus.Passed {
			//	logs.WithF(h.fields).Error("Smart status is failed")
			//}

			//h.storeInfo()


		case <- h.stop:
			ticker.Stop()
			return
		}
	}
}
