package agent

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

// store disk info in db
type HandlerDb struct {
	CommonHandler
}

func (h *HandlerDb) Start() {
	h.storeInfo()
}

func (h *HandlerDb) Stop() {

}

///////////////////////////////////

func (h *HandlerDb) storeInfo() {
	disk, err := h.server.ScanDisk(h.path)
	if err != nil {
		logs.WithE(err).Error("Failed to scan disk")
	}

	if err := hdm.HDM.DBDisk().Save(disk); err != nil {
		logs.WithE(err).Error("Failed to save disk in db")
	}
}
