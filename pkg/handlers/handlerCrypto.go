package handlers

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
)

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
	passwordSet := h.manager.PassService.Watch()
	defer h.manager.PassService.Unwatch(passwordSet)

	if err := h.open(); err != nil {
		logs.WithE(err).Error("Failed to open crypto")
	}

	for {
		select {
		case <-passwordSet:
			if err := h.open(); err != nil {
				logs.WithE(err).Error("Failed to open crypto")
			}
		case <-h.stop:
			return
		}
	}
}

func (h *HandlerCrypto) open() error {

	if !h.manager.PassService.IsSet() {
		logs.WithF(h.fields).Debug("Password is not set, cannot open")
		return nil
	}

	buffer, err := h.manager.PassService.Get()
	if err != nil {
		return errs.WithEF(err, h.fields, "Failed to get password from password service")
	}
	defer buffer.Destroy()

	if err := h.manager.BlockDevice.LuksOpen(buffer); err != nil {
		return errs.WithEF(err, h.fields, "Failed to Open")
	}
	return nil
}