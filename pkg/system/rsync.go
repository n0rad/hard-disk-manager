package system

import (
	"github.com/alessio/shellescape"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"path"
)

type Rsync struct {
	SourceExec             runner.Exec
	SourceMountPoint       string
	SourceInFilesystemPath string
	//SourceFilesystem       BlockDevice

	TargetExec             runner.Exec
	TargetMountPoint       string
	TargetInFilesystemPath string
	//TargetFilesystem       BlockDevice

	Delete bool

	sourceFullPath string
	targetFullPath string

	fields data.Fields
}

func (r *Rsync) Init() error {
	if r.SourceMountPoint == "" {
		return errs.With("Source filesystem cannot be empty or is not mounted")
	}
	if r.TargetMountPoint == "" {
		return errs.With("Target filesystem cannot be empty or is not mounted")
	}
	if r.SourceInFilesystemPath == "" {
		return errs.With("Source in filesystem path cannot be empty")
	}
	if r.TargetInFilesystemPath == "" {
		return errs.With("Target in filesystem path cannot be empty")
	}

	if r.TargetInFilesystemPath[0] != '/' {
		r.TargetInFilesystemPath = "/" + r.TargetInFilesystemPath
	}

	if r.SourceInFilesystemPath[0] != '/' {
		r.SourceInFilesystemPath = "/" + r.SourceInFilesystemPath
	}

	r.targetFullPath = r.TargetMountPoint + r.TargetInFilesystemPath
	r.sourceFullPath = r.SourceMountPoint + r.SourceInFilesystemPath

	r.fields = data.WithField("source", r.sourceFullPath).WithField("target", r.targetFullPath)
	return nil
}

func (r *Rsync) SourceSize() (int, error) {
	return DirectorySize(r.sourceFullPath, r.SourceExec)
}

func (r *Rsync) TargetSize() (int, error) {
	targetPath := r.targetFullPath + "/" + path.Base(r.sourceFullPath)
	return DirectorySize(targetPath, r.TargetExec)
}

func (r *Rsync) Rsyncable() (error, error) {
	if r.TargetMountPoint == "" {
		return errs.WithF(r.fields.WithField("disk", r.TargetMountPoint), "Disk is not mounted"), nil
	}

	//sourceSize, err := r.SourceSize()
	//if err != nil {
	//	return nil, errs.WithEF(err, r.fields, "Cannot get directory size")
	//}
	//
	//targetSize, err := r.TargetSize()
	//if err != nil {
	//	return nil, errs.WithEF(err, r.fields, "Cannot get directory size")
	//}
	//
	//targetAvailable, err := r.TargetFilesystem.SpaceAvailable()
	//if err != nil {
	//	return nil, errs.WithEF(err, r.fields, "Cannot get TargetLabel available space")
	//}
	//
	//if sourceSize > targetSize+targetAvailable {
	//	return errs.WithF(data.WithField("sourceSize", sourceSize).
	//		WithField("targetSize", targetSize).
	//		WithField("targetAvailable", targetAvailable), "Not enough space to backup"), nil
	//}
	return nil, nil
}

// TODO support sync to other server
func (r *Rsync) RSync() error {
	logs.WithFields(r.fields).Debug("rsync")
	why, err := r.Rsyncable()
	if err != nil {
		return errs.WithEF(err, r.fields, "Failed to see if directory is backupable")
	}
	if why != nil {
		logs.WithEF(why, r.fields).Warn("Directory is not backupable")
		return nil
	}

	if std, err := r.TargetExec.ExecGetStd("mkdir", "-p ", shellescape.Quote(r.targetFullPath)); err != nil {
		return errs.WithEF(err, r.fields.WithField("path", r.targetFullPath).WithField("std", std), "Failed to create target backup path")
	}

	args := []string{"-avP", "--itemize-changes", r.sourceFullPath, r.targetFullPath}
	if r.Delete {
		args = append(args, "--delete")
	}

	logs.WithFields(r.fields).Info("Running backup")
	_, stderr, err := r.SourceExec.ExecGetStdoutStderr("rsync", args...)
	if err != nil {
		return errs.WithEF(err, r.fields.WithField("stderr", stderr), "Backup failed")
	}
	return nil
}
