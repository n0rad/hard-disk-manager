package trash

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
)

func (h *Config) FillFromFile(filesystem system.BlockDeviceOLD, file string) error {
	bytes, err := filesystem.ExecShell("cat " + shellescape.Quote(file))
	if err != nil {
		return errs.WithEF(err, data.WithField("file", file), "Failed to cat file")
	}

	if err := yaml.Unmarshal([]byte(bytes), h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to parse hdm file")
	}

	if err := h.Init(file); err != nil {
		return errs.WithEF(err, h.fields.WithField("content", string(bytes)), "Failed to init hdm file")
	}

	return nil
}
