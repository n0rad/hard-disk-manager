package system

import (
	"github.com/n0rad/go-erlog/errs"
)

func (b BlockDevice) Format(partType string, label string) error {
	if b.HasChildren() {
		return errs.WithF(b.fields, "Cannot format, has children")
	}

	if std, err := b.exec.ExecGetStd("mkfs." + partType, "-L", label, "-f", b.Path); err != nil {
		return errs.WithEF(err, b.fields.WithField("std", std), "Failed to make filesystem")
	}
	return nil
}