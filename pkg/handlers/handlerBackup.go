package handlers

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
