package handlers

import (
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
	//TODO init elsewhere
	systemd := system.Systemd{}
	systemd.Init(h.manager.BlockDevice.GetExec())
	//mountPath, err := systemd.SystemdMountPath(h.manager.BlockDevice.Path)

	//if err != nil {
	//	logs.WithEF(err, data.WithField("path", h.manager.BlockDevice.Path)).Trace("Failed to get systemd mount path")
		mountPath := h.DefaultMountPath + "/" + h.manager.BlockDevice.Name
		if h.manager.BlockDevice.Label != "" {
			mountPath = h.DefaultMountPath + "/" + h.manager.BlockDevice.Label
		}
	//}

	if err := h.manager.BlockDevice.Mount(mountPath); err != nil {
		if err := h.manager.BlockDevice.Umount(mountPath); err != nil {
			logs.WithE(err).Warn("Failed to cleanup failed mount")
		}
		return err
	}
	return nil
}