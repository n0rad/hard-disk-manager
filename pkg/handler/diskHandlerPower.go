package handler

func init() {
	diskHandlerBuilders["disk-power"] = diskHandlerBuilder{
		new: func() DiskHandler {
			return &HandlerPower{}
		},
	}
}

type HandlerPower struct {
	CommonDiskHandler
}

func (h *HandlerPower) Remove() error {
	return h.manager.block.PutInSleepNow()
}
