package system

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"strconv"
	"strings"
)

func DirectorySize(path string, exec runner.Exec) (int, error) {
	_, stderr, err := exec.ExecGetStdoutStderr("test", "-d", path)
	if err != nil {
		return 0, errs.WithEF(err, data.WithField("path", path), "Failed to check directory exists")
	}

	stdout, stderr, err := exec.ExecGetStdoutStderr("du", "-s", path)
	if err != nil {
		return 0, errs.WithEF(err, data.WithField("path", path).WithField("stderr", stderr), "Failed to get directory size")
	}
	split := strings.Split(stdout, "\t")

	size, err := strconv.Atoi(strings.TrimSpace(split[0]))
	if err != nil {
		return 0, errs.WithEF(err, data.WithField("path", path), "Failed to parse 'du' result")
	}
	return size, nil
}
