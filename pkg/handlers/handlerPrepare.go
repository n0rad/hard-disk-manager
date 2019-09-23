package handlers

import "github.com/n0rad/go-erlog/logs"

// prepare disk if no partitions
type HandlerPrepare struct {
	CommonHandler
}

func (h *HandlerPrepare) Start() {
	//disk, err := h.server.ScanDisk(h.path)
	//if err != nil {
	//	logs.WithE(err).Error("Failed to scan disk")
	//	return
	//}

	if len(h.disk.Children) == 0 {
		logs.WithF(h.fields).Warn("Disk has not partitions and should be prepared")
	}
}
