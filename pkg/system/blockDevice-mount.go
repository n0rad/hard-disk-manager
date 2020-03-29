package system

import (
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
)

func (b *BlockDevice) IsMounted() bool {
	return b.Mountpoint != ""
}

func (b *BlockDevice) Mount(mountPath string) error {
	logs.WithFields(b.fields.WithField("mountPath", mountPath)).Debug("Mount")
	if mountPath == "" {
		return errs.WithF(b.fields, "mountPath cannot be empty")
	}

	// TODO cannot differentiate exec fail from not mounted
	if _, err := b.exec.ExecShellGetStd("cat /proc/mounts | cut -f1,2 -d' ' | grep '" + b.Path + " " + mountPath + "$'"); err == nil {
		logs.WithF(b.fields).Debug("Directory is already mounted")
		return nil
	}

	if std, err := b.exec.ExecGetStd("mkdir", "-p", mountPath); err != nil {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("std", std), "Failed to create mount directory")
	}

	out, err := b.exec.ExecGetStdout("ls", "-A", mountPath)
	if err != nil {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", out), "Failed to ls on mount path")
	}
	if string(out) != "" {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", out), "Directory is not empty")
	}

	if std, err := b.exec.ExecShellGetStd("! cat /proc/mounts | cut -f2 -d' ' | grep " + mountPath + "$"); err != nil {
		logs.WithEF(err, b.fields.WithField("path", mountPath).WithField("std", std)).Trace("Already mounted")
		return nil
	}

	if std, err := b.exec.ExecGetStd("mount", b.Path, mountPath); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", std).WithField("target", mountPath), "Failed to mount")
	}
	return nil
}

// mountPath is required to cleanup directory in case it's not mounted, but directory was created
func (b *BlockDevice) Umount(mountPath string) error {
	logs.WithFields(b.fields).Debug("Umount")
	if mountPath == "" {
		mountPath = b.Mountpoint
	}

	if mountPath != "" {
		if std, err := b.exec.ExecGetStd("umount", mountPath); err != nil {
			return errs.WithEF(err, b.fields.WithField("std", std), "Failed to unmount")
		}
	}

	if std, err := b.exec.ExecGetStd("rmdir", mountPath); err != nil {
		logs.WithEF(err, b.fields.WithField("std", std).WithField("path", mountPath)).Warn("Failed to cleanup mount path")
	}

	return nil
}
