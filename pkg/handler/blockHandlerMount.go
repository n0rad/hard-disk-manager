package handler

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
	BlockHandlers[handlerNameMount] = BlockHandlerBuilder{
		new: func() BlockHandler {
			return &HandlerMount{}
		},
	}
}

type HandlerMount struct {
	CommonBlockHandler
	DefaultMountPath string
}

func (h *HandlerMount) Init(name string, manager *BlockManager) {
	h.CommonBlockHandler.Init(name, manager)

	h.DefaultMountPath = manager.GetHDM().DefaultMountPath
}

func (h *HandlerMount) Add() error {
	if !utils.SliceContains(system.Filesystems, h.manager.block.Fstype) {
		logs.WithF(h.GetFields().WithField("fstype", h.manager.block.Fstype)).Trace("Unsupported fstype")
		return nil
	}

	logs.WithFields(h.GetFields()).Info("Add")

	if h.manager.block.Mountpoint != "" {
		logs.WithF(h.GetFields()).Debug("Already mounted")
		return nil
	}

	mountPath := h.DefaultMountPath + "/" + h.manager.block.GetUsableLabel()
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to create mount directory")
	}

	if err := h.manager.block.Mount(mountPath); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to mount")
	}

	return nil
}

func (h *HandlerMount) Remove() error {
	logs.WithFields(h.GetFields()).Info("Remove")

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
			return errs.WithEF(err, h.GetFields().WithField("dir", h.manager.block.Mountpoint), "Failed to remove mount dir")
		}
	}

	return nil
}
