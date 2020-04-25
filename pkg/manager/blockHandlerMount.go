package manager

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
	blockHandlers[handlerNameMount] = blockHandlerBuilder{
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

	h.DefaultMountPath = manager.hdm.DefaultMountPath
}

func (h *HandlerMount) Add() error {
	if !utils.SliceContains(system.Filesystems, h.manager.block.Fstype) {
		logs.WithF(h.GetFields().WithField("fstype", h.manager.block.Fstype)).Trace("Unsupported fstype")
		return nil
	}

	logs.WithFields(h.GetFields()).Info("Add")

	if h.manager.block.Mountpoint != "" {
		logs.WithF(h.GetFields()).Debug("Already mounted")

		//TODO -------
		return h.registerFsManager()
	}

	mountPath := h.DefaultMountPath + "/" + h.manager.block.GetUsableLabel()
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to create mount directory")
	}

	if err := h.manager.block.Mount(mountPath); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to mount")
	}

	//TODO -------
	return h.registerFsManager()
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

func (h *HandlerMount) registerFsManager() error {
	manager := &FsManager{}
	if err := manager.Init(h.manager, h.manager.block); err != nil {
		return errs.WithEF(err, h.fields, "Failed to init fsManager")
	}

	h.manager.children[h.manager.block.Mountpoint] = manager
	return nil
}