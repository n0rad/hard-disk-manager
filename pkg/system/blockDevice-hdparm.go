package system

import (
	"github.com/n0rad/go-erlog/errs"
)

func (b *BlockDevice) SecureErase() error {
	if b.HasChildren() {
		return errs.WithF(b.fields, "Cannot erase, has children")
	}

	if std, err := b.exec.ExecGetStd("hdparm", "--user-master", "u", "--security-erase", "Nine", b.Path); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", std), "Fail to erase disk")
	}
	return nil
}

func (b *BlockDevice) PutInSleepNow() error {
	if std, err := b.exec.ExecGetStd("hdparm", "-y", b.Path); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", std), "Failed to put disk in sleep")
	}
	return nil
}
