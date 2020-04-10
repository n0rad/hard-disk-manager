package handler

func init() {
	DiskHandlerBuilders["power"] = diskHandlerBuilder{
		//block.HandlerFilter{Type: "disk"},
		New: func() DiskHandler {
			return &HandlerPower{
				CommonDiskHandler: CommonDiskHandler{
					CommonBlockHandler: CommonBlockHandler{
						CommonHandler: CommonHandler{
							HandlerName: "power",
						},
					},
				},
			}
		},
	}
}

type HandlerPower struct {
	CommonDiskHandler
}

func (h *HandlerPower) Remove() error {
	return h.manager.Block.PutInSleepNow()
}
