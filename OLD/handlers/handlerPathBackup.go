package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "path"},
		func() Handler {
			return &HandlerBackup{
				CommonHandler: CommonHandler{
					handlerName: "handlerBackup",
				},
			}
		},
	})
}

type HandlerBackup struct {
	CommonHandler
}

//func (h *HandlerBackup) Init(manager *BlockManager) {
//	h.CommonHandler.Init(manager)
//}

func (h *HandlerBackup) Start() {



	//func (h *Config) RunBackups(disks system.Disks) error {
	//	for _, backup := range h.Backups {
	//		if err := backup.Backup(disks); err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//}

	for _, v := range h.manager.config.Backups {
		b := hdm.Backup{
			BackupConfig: v,
		}
		if err := b.Init(h.manager.config.GetConfigPath(), h.manager.BlockDevice, hdm.HDM.Servers); err != nil {
			logs.WithE(err).Error("Failed to init backup")
			continue
		}

		if err := b.Backup(); err !=nil {
			logs.WithE(err).Error("Failed to backup")
		}
	}

	// TODO
	// check latest backup date



	//logs.WithFields(h.fields).Warn("handle backup")
	<-h.stop
	//ticker := time.NewTicker(h.CheckInterval)
	//for {
	//	select {
	//	case <- ticker.C:
	//		logs.WithFields(h.fields).Debug("Time to check disk status")
	//	case <- h.stop:
	//		ticker.Stop()
	//		return
	//	}
	//}
}
