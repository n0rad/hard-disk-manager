package agent

// store disk info in db
type HandlerHealthCheck struct {
	CommonHandler
}

func (h *HandlerHealthCheck) Start() {

	//err := h.disk.FillFromSmartctl()
	//
	//
	//if !h.disk.SmartResult.SmartStatus.Passed {
	//	logs.WithF(h.fields).Error("Smart status is failed")
	//}

	//h.storeInfo()

}
