package block

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
	"os"
	"syscall"
)

const handlerNameMount = "mount"

func init() {
	handlers[handlerNameMount] = handler{
		HandlerFilter{FSType: ""},
		func() Handler {
			return &HandlerMount{
				CommonHandler: CommonHandler{
					handlerName: handlerNameMount,
				},
			}
		},
	}
}

type HandlerMount struct {
	CommonHandler
	DefaultMountPath string
}

func (h *HandlerMount) Init(manager *Manager) {
	h.CommonHandler.Init(manager)
	h.DefaultMountPath = manager.hdm.DefaultMountPath
}

func (h *HandlerMount) Add() error {
	if !utils.SliceContains(system.Filesystems, h.manager.block.Fstype) {
		logs.WithF(h.fields.WithField("fstype", h.manager.block.Fstype)).Trace("Unsupported fstype")
		return nil
	}

	logs.WithFields(h.fields).Info("Add")

	if h.manager.block.Mountpoint != "" {
		logs.WithF(h.fields).Debug("Already mounted")
		return nil
	}

	mountPath := h.DefaultMountPath + "/" + h.manager.block.GetUsableLabel()
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return errs.WithEF(err, h.fields, "Failed to create mount directory")
	}

	if err := h.manager.block.Mount(mountPath); err != nil {
		return errs.WithEF(err, h.fields, "Failed to mount")
	}

	return nil
}

func (h *HandlerMount) Remove() error {
	logs.WithFields(h.fields).Info("Remove")

	// TODO free
	// kill lsof

	if h.manager.block.Mountpoint != "" {
		if err := h.manager.block.Umount(h.manager.block.Mountpoint); err != nil {
			return err
		}
	}

	stat, err := os.Stat(h.manager.block.Mountpoint)
	if err == nil && stat.IsDir() {
		if err := syscall.Rmdir(h.manager.block.Mountpoint); err != nil {
			return errs.WithEF(err, h.fields.WithField("dir", h.manager.block.Mountpoint), "Failed to remove mount dir")
		}
	}

	return nil
}
