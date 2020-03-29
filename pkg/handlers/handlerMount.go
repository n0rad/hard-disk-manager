package handlers

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{FSType:""},
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
	if !utils.SliceContains(system.Filesystems, h.manager.BlockDevice.Fstype) {
		logs.WithField("fstype", h.manager.BlockDevice.Fstype).Debug("Unsupported fstype")
		return nil
	}


	//TODO init elsewhere
	systemd := system.Systemd{}
	systemd.Init(h.manager.BlockDevice.GetExec())
	//mountPath, err := systemd.SystemdMountPath(h.manager.BlockDevice.Path)

	//if err != nil {
	//	logs.WithEF(err, data.WithField("path", h.manager.BlockDevice.Path)).Trace("Failed to get systemd mount path")

	mountPath := h.DefaultMountPath + "/" + h.manager.BlockDevice.GetUsableLabel()

	if h.manager.BlockDevice.Mountpoint != "" {
		logs.WithF(h.manager.BlockDevice.GetFields()).Debug("Already mounted")
		return nil
	}


	if err := h.manager.BlockDevice.Mount(mountPath); err != nil {
		logs.WithE(err).Warn("Failed to mount")
		if err := h.manager.BlockDevice.Umount(mountPath); err != nil {
			logs.WithE(err).Warn("Failed to cleanup failed mount")
		}
		return err
	}
	return nil
}