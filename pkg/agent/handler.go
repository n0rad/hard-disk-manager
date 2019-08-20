package agent


type Handler interface {
	Init(path string)
	Start()
	Stop()
	//Handle(event)
}


type CommonHandler struct {
	path string
}

func (h *CommonHandler) Init(path string) {
	h.path = path
}



type Event struct {

}
