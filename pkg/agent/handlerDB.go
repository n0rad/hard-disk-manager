package agent

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"time"
)

// store disk info in db
type HandlerDb struct {
	CommonHandler
	storeInterval time.Duration
}

func (h *HandlerDb) Init(manager *DiskManager) {
	h.CommonHandler.Init(manager)

	if h.storeInterval == 0 {
		h.storeInterval = 10 * time.Second
	}
}

func (h *HandlerDb) Start() {
	h.storeInfo()

	// TODO handle partitions changes
	ticker := time.NewTicker(h.storeInterval)
	for {
		select {
		case <- ticker.C:
			logs.WithFields(h.fields).Debug("Time to store info")
			h.storeInfo()
		case <- h.stop:
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
