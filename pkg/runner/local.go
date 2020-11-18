package runner

import (
	"bytes"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"io"
	"os"
	"os/exec"
	"strings"
)

var Local Exec = LocalExec{}

type LocalExec struct {
	//UnSudo bool
}

func (s LocalExec) String() string {
	return "local"
}

func (s LocalExec) Close() {}

func (s LocalExec) Exec(head string, args ...string) error {
	return s.ExecStdinStdoutStderr(os.Stdin, os.Stdout, os.Stderr, head, args...)
}

func (s LocalExec) ExecGetStd(head string, args ...string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr(head, args...)
	stdout += stderr
	return stdout, err
}

func (s LocalExec) ExecGetStdout(head string, args ...string) (string, error) {
	stdout, _, err := s.ExecGetStdoutStderr(head, args...)
	return stdout, err
}

func (s LocalExec) ExecGetStdoutStderr(head string, args ...string) (string, string, error) {
	return s.ExecSetStdinGetStdoutStderr(nil, head, args...)
}

func (s LocalExec) ExecSetStdinGetStdoutStderr(stdin io.Reader, head string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := s.ExecStdinStdoutStderr(stdin, &stdout, &stderr, head, args...)
	return strings.TrimSpace(stdout.String()), stderr.String(), err
}

func (s LocalExec) ExecStdinStdoutStderr(stdin io.Reader, stdout io.Writer, stderr io.Writer, head string, args ...string) error {
	commandDebug := strings.Join([]string{head, " ", strings.Join(args, " ")}, " ")
	if logs.IsDebugEnabled() {
		logs.WithField("command", commandDebug).Debug("Running command")
	}
	cmd := exec.Command(head, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return errs.WithEF(err, data.WithField("command", commandDebug), "Failed to start command")
	}
	return cmd.Wait()
}

/////

func (s LocalExec) ExecShellGetStd(cmd string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr("bash", "-o", "pipefail", "-c", cmd)
	stdout += stderr
	return stdout, err
}

func (s LocalExec) ExecShellGetStdout(cmd string) (string, error) {
	return s.ExecGetStdout("bash", "-o", "pipefail", "-c", cmd)
}
