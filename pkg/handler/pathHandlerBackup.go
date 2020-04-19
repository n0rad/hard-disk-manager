package handler

func init() {
	pathHandlerBuilders["backup"] = pathHandlerBuilder{
		new: func() PathHandler {
			return &PathHandlerBackup{}
		},
	}
}

type PathHandlerBackup struct {
	CommonPathHandler
}

func (h *PathHandlerBackup) Add() error {

	//for _, v := range h.manager.config.Backups {
	//	b := hdm.Backup{
	//		BackupConfig: v,
	//	}
	//	if err := b.Init(h.manager.config.GetConfigPath(), h.manager.BlockDevice, hdm.HDM.Servers); err != nil {
	//		logs.WithE(err).Error("Failed to init backup")
	//		continue
	//	}
	//
	//	if err := b.Backup(); err !=nil {
	//		logs.WithE(err).Error("Failed to backup")
	//	}
	//}

	return nil
}
