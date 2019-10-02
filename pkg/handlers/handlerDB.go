package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"time"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "disk"},
		func() Handler {
			return &HandlerDb{
				CommonHandler: CommonHandler{
					handlerName: "db",
				},
			}
		},
	})
}

// store disk info in db
type HandlerDb struct {
	CommonHandler
	StoreInterval time.Duration
}

func (h *HandlerDb) Init(manager *BlockManager) {
	h.CommonHandler.Init(manager)

	if h.StoreInterval == 0 {
		h.StoreInterval = 10 * time.Second
	}
}

func (h *HandlerDb) Start() {
	h.storeInfo()

	// TODO handle partitions changes
	ticker := time.NewTicker(h.StoreInterval)
	for {
		select {
		case <-ticker.C:
			logs.WithFields(h.fields).Debug("Time to store info")
			h.storeInfo()
		case <-h.stop:
			ticker.Stop()
			return
		}
	}
}

///////////////////////////////////

func (h *HandlerDb) storeInfo() {
	disk, err := h.server.ScanDisk(h.manager.Path)
	if err != nil {
		logs.WithE(err).Error("Failed to scan disk")
		return
	}

	if err := hdm.HDM.DBDisk().Save(disk); err != nil {
		logs.WithE(err).Error("Failed to save disk in db")
	}
}
