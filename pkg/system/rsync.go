package system

import (
	"github.com/alessio/shellescape"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"path"
	"strconv"
	"strings"
)

type Rsync struct {
	SourceFilesystem       BlockDeviceOLD
	SourceInFilesystemPath string

	TargetFilesystem       BlockDeviceOLD
	TargetInFilesystemPath string

	Delete bool

	sourceFullPath string
	targetFullPath string

	fields data.Fields
}

func (r *Rsync) Init() error {
	if r.SourceFilesystem.Mountpoint == "" {
		return errs.With("Source filesystem cannot be empty or is not mounted")
	}
	if r.TargetFilesystem.Mountpoint == "" {
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

	r.targetFullPath = r.TargetFilesystem.Mountpoint + r.TargetInFilesystemPath
	r.sourceFullPath = r.SourceFilesystem.Mountpoint + r.SourceInFilesystemPath

	r.fields = data.WithField("source", r.sourceFullPath).WithField("target", r.targetFullPath)
	return nil
}

func (r *Rsync) SourceSize() (int, error) {
	out, err := r.SourceFilesystem.ExecShell("du -s " + shellescape.Quote(r.sourceFullPath) + " | cut -f1")
	if err != nil {
		return 0, errs.WithEF(err, r.fields, "Failed to get directory size")
	}
	size, err := strconv.Atoi(out)
	if err != nil {
		return 0, errs.WithEF(err, r.fields, "Failed to parse 'du' result")
	}
	return size, nil
}

func (r *Rsync) TargetSize() (int, error) {
	targetPath := r.targetFullPath + "/" + path.Base(r.sourceFullPath)
	_, err := r.TargetFilesystem.ExecShell("sudo test -d " + shellescape.Quote(targetPath))
	if err != nil {
		return 0, nil
	}

	bytes, err := r.TargetFilesystem.ExecShell("sudo du -s " + shellescape.Quote(targetPath) + " | cut -f1")
	if err != nil {
		return 0, errs.WithEF(err, r.fields, "Failed to get directory size")
	}
	size, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		return 0, errs.WithEF(err, r.fields, "Failed to parse 'du' result")
	}
	return size, nil
}

func (r *Rsync) Rsyncable() (error, error) {
	if r.TargetFilesystem.Mountpoint == "" {
		return errs.WithF(r.fields.WithField("disk", r.TargetFilesystem), "Disk is not mounted"), nil
	}

	sourceSize, err := r.SourceSize()
	if err != nil {
		return nil, errs.WithEF(err, r.fields, "Cannot get directory size")
	}

	targetSize, err := r.TargetSize()
	if err != nil {
		return nil, errs.WithEF(err, r.fields, "Cannot get directory size")
	}

	targetAvailable, err := r.TargetFilesystem.SpaceAvailable()
	if err != nil {
		return nil, errs.WithEF(err, r.fields, "Cannot get TargetLabel available space")
	}

	if sourceSize > targetSize+targetAvailable {
		return errs.WithF(data.WithField("sourceSize", sourceSize).
			WithField("targetSize", targetSize).
			WithField("targetAvailable", targetAvailable), "Not enough space to backup"), nil
	}
	return nil, nil
}

func (r *Rsync) RSync() error {
	why, err := r.Rsyncable()
	if err != nil {
		return errs.WithEF(err, r.fields, "Failed to see if directory is backupable")
	}
	if why != nil {
		logs.WithEF(why, r.fields).Warn("Directory is not backupable")
		return nil
	}

	if _, err := r.TargetFilesystem.ExecShell("sudo", "mkdir", "-p ", shellescape.Quote(r.targetFullPath)); err != nil {
		return errs.WithEF(err, r.fields.WithField("path", r.targetFullPath), "Failed to create target backup path")
	}

	deleteIfSourceRemoved := ""
	if r.Delete {
		deleteIfSourceRemoved = "--delete"
	}

	logs.WithFields(r.fields).Info("Running backup")
	_, err = r.SourceFilesystem.ExecShell("sudo Rsync -avP " + deleteIfSourceRemoved + " --itemize-changes " + shellescape.Quote(r.sourceFullPath) + " " + shellescape.Quote(r.targetFullPath)) // TODO support sync to other server
	if err != nil {
		return errs.WithEF(err, r.fields, "Backup failed")
	}
	return nil
}
