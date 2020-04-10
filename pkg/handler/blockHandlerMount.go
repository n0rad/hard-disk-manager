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
		New: func() BlockHandler {
			return &HandlerMount{
				CommonBlockHandler: CommonBlockHandler{
					CommonHandler: CommonHandler{
						HandlerName: "mount",
					},
				},
			}
		},
	}
}

type HandlerMount struct {
	CommonBlockHandler
	DefaultMountPath string
}

func (h *HandlerMount) Init(manager *BlockManager) {
	h.CommonBlockHandler.Init(manager)

	h.DefaultMountPath = manager.GetHDM().DefaultMountPath
}

func (h *HandlerMount) Add() error {
	if !utils.SliceContains(system.Filesystems, h.manager.Block.Fstype) {
		logs.WithF(h.GetFields().WithField("fstype", h.manager.Block.Fstype)).Trace("Unsupported fstype")
		return nil
	}

	logs.WithFields(h.GetFields()).Info("Add")

	if h.manager.Block.Mountpoint != "" {
		logs.WithF(h.GetFields()).Debug("Already mounted")
		return nil
	}

	mountPath := h.DefaultMountPath + "/" + h.manager.Block.GetUsableLabel()
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to create mount directory")
	}

	if err := h.manager.Block.Mount(mountPath); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to mount")
	}

	return nil
}

func (h *HandlerMount) Remove() error {
	logs.WithFields(h.GetFields()).Info("Remove")

	// TODO free
	// kill lsof

	if h.manager.Block.Mountpoint != "" {
		if err := h.manager.Block.Umount(h.manager.Block.Mountpoint); err != nil {
			return err
		}
	}

	stat, err := os.Stat(h.manager.Block.Mountpoint)
	if err == nil && stat.IsDir() {
		if err := syscall.Rmdir(h.manager.Block.Mountpoint); err != nil {
			return errs.WithEF(err, h.GetFields().WithField("dir", h.manager.Block.Mountpoint), "Failed to remove mount dir")
		}
	}

	return nil
}
