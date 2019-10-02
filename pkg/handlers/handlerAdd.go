package handlers

type HandlerAdd struct {
	CommonHandler
}

func (h *HandlerAdd) Start() {

	password := h.manager.PassService.Watch()
	<-password

	//buffer, err := h.manager.PassService.Get()
	//if err != nil {
	//	logs.WithEF(err, h.fields).Error("Cannot get password to add disk")
	//	return
	//}

	//if err := h.disk.AddBlockDevice(buffer); err != nil {
	//	logs.WithEF(err, h.fields).Error("Failed to add disk")
	//}



	//disk, err := h.server.ScanDisk(h.path)
	//if err != nil {
	//	logs.WithE(err).Error("Failed to scan disk")
	//	return
	//}

	//h.disk.AddBlockDevice()
}

