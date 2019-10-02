package handlers

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{FSType: "crypto_LUKS"},
		func() Handler {
			return &HandlerCrypto{
				CommonHandler: CommonHandler{
					handlerName: "handlerCrypto",
				},
			}
		},
	})
}

type HandlerCrypto struct {
	CommonHandler
}

func (h *HandlerCrypto) Start() {
	//disk, err := h.server.ScanDisk(h.manager.Path)


	// start or event of password change

}