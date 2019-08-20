package agent

// store disk info in db
type HandlerDb struct {
	CommonHandler
}

func (h *HandlerDb) Start() {
	//disk, err := server.ScanDisk(path)
	//if err != nil {
	//	return DiskManager{}, errs.WithEF(err, data.WithField("path", path), "Failed to scan disk")
	//}

	// store disk in db
}

func (h *HandlerDb) Stop() {

}
