package system

import (
	"github.com/awnumar/memguard"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"os"
)

func (b *BlockDevice) IsLuksNameUsed() (bool, error) {
	if _, err := os.Stat("/dev/mapper/" + b.GetUsableLabel()); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// check if luks is still open for current block
// this can happen when the device is removed / add back.
// device mapper is till there, but its not linked to the hosting block device
func (b *BlockDevice) IsLuksUsed() bool {
	if _, err := b.exec.ExecShellGetStd("dmsetup table "+b.GetUsableLabel()+" | grep $(cryptsetup luksUUID "+b.Path+")"); err != nil {
		return false // TODO difference nil from error
	}
	return true
}

// sudo cryptsetup luksUUID /dev/sdc1
// ls /dev/disk/by-id/dm-uuid-CRYPT-LUKS2-20a8278d326f41bc84385f95d422f1de-yopla
func (b *BlockDevice) IsLuksOpen() bool {
	if b.Fstype != "crypto_LUKS" {
		return false
	}
	if b.HasChildren() {
		return true
	}
	return false
}

func (b *BlockDevice) LuksOpen(cryptPassword *memguard.LockedBuffer) error {
	logs.WithFields(b.fields).Debug("Disk luksOpen")
	if b.Fstype != "crypto_LUKS" {
		return errs.WithF(b.fields.WithField("fstype", b.Fstype), "Cannot luks open, not a crypto block device")
	}

	if b.HasChildren() {
		logs.WithFields(b.fields).Warn("Opening an already open crypto")
		return nil
	}

	volumeName := b.Partlabel
	if volumeName == "" {
		volumeName = b.Name
	}

	if stdout, stderr, err := b.exec.ExecSetStdinGetStdoutStderr(cryptPassword.Reader(), "cryptsetup",
		"luksOpen",
		b.Path,
		volumeName,
		"-"); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", stdout + stderr), "Failed to open luks")
	}

	return nil
}

func (b *BlockDevice) LuksFormat(cryptPassword *memguard.LockedBuffer) error {
	logs.WithFields(b.fields).Debug("Luks format")
	if b.Parttype != luksPartitionCode {
		return errs.WithF(b.fields.WithField("parttype", b.Parttype), "Cannot luks format, not a luks partition")
	}

	if b.HasChildren() {
		return errs.WithF(b.fields, "Cannot luks format, has children")
	}

	if stdout, stderr, err := b.GetExec().ExecSetStdinGetStdoutStderr(cryptPassword.Reader(), "cryptsetup",
		"--verbose",
		"--hash=sha512",
		"--cipher=aes-xts-benbi:sha512",
		"--key-size=512",
		"luksFormat",
		b.Path,
		"-"); err != nil {
		return errs.WithEF(err,b.fields.WithField("std", stdout + stderr), "Fail to crypt partition")
	}
	return nil
}

func (b *BlockDevice) LuksClose() error {
	logs.WithFields(b.fields).Debug("Luks close")
	if std, err := b.exec.ExecGetStd("cryptsetup", "luksClose", b.Path); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", std), "Failed to luks close")
	}
	return nil
}
