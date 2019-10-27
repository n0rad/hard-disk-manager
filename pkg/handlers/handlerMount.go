package handlers

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{},
		func() Handler {
			return &HandlerMount{
				CommonHandler: CommonHandler{
					handlerName: "mount",
				},
			}
		},
	})
}

type HandlerMount struct {
	CommonHandler
	DefaultMountPath string
}

func (h *HandlerMount) Init(manager *BlockManager) {
	h.CommonHandler.Init(manager)
	h.DefaultMountPath = hdm.HDM.DefaultMountPath // TODO
}

func (h *HandlerMount) Start() {

	if err := h.tryMount(); err != nil {
		logs.WithEF(err, h.fields).Debug("Failed to mount")
	}

	<-h.stop
}

func (h *HandlerMount) tryMount() error {
	b, err := h.server.GetBlockDevice(h.manager.Path)
	if err != nil {
		return errs.WithEF(err, h.fields, "Failed to get blockDevice")
	}

	mountPath, err := system.SystemdMountPath(b.Path)
	if err != nil {
		logs.WithEF(err, data.WithField("path", b.Path)).Trace("Failed to get systemd mount path")
		mountPath = h.DefaultMountPath + "/" + b.Name
		if b.Label != "" {
			mountPath = h.DefaultMountPath + "/" + b.Label
		}
	}

	if err := b.Mount(mountPath); err != nil {
		if err := b.Umount(mountPath); err != nil {
			logs.WithE(err).Warn("Failed to cleanup failed mount")
		}
		return err
	}
	return nil
}