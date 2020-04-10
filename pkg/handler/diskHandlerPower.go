package handler

func init() {
	DiskHandlerBuilders["disk-power"] = diskHandlerBuilder{
		New: func() DiskHandler {
			return &HandlerPower{}
		},
	}
}

type HandlerPower struct {
	CommonDiskHandler
}

func (h *HandlerPower) Remove() error {
	return h.manager.Block.PutInSleepNow()
}
