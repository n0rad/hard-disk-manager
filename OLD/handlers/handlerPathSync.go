package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

// TODO on inotify on a file
// on periodic schedule

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "path"},
		func() Handler {
			return &HandlerSync{
				CommonHandler: CommonHandler{
					handlerName: "handlerSync",
				},
			}
		},
	})
}

type HandlerSync struct {
	CommonHandler
}

func (h *HandlerSync) Start() {
	for _, v := range h.manager.config.Syncs {
		s := hdm.Sync{
			SyncConfig: v,
		}
		if err := s.Init(h.manager.config.GetConfigPath(), h.manager.BlockDevice, hdm.HDM.Servers); err != nil {
			logs.WithE(err).Error("Failed to init sync")
			continue
		}

		if err := s.Sync(); err !=nil {
			logs.WithE(err).Error("Failed to sync")
		}
	}

	<-h.stop
}
