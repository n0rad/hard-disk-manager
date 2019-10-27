package runner

import (
	"bytes"
	"github.com/n0rad/go-erlog/logs"
	"os/exec"
	"strings"
)

var Local Runner = LocalRunner{}

type LocalRunner struct {
	//UnSudo bool
}

func (s LocalRunner) ExecShellGetStd(cmd string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr("bash", "-o", "pipefail", "-c", cmd)
	stdout += stderr
	return stdout, err
}

func (s LocalRunner) ExecShellGetStdout(cmd string) (string, error) {
	return s.ExecGetStdout("bash", "-o", "pipefail", "-c", cmd)
}

//

func (s LocalRunner) ExecGetStd(head string, args ...string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr(head, args...)
	stdout += stderr
	return stdout, err
}

func (s LocalRunner) ExecGetStdout(head string, args ...string) (string, error) {
	stdout, _, err := s.ExecGetStdoutStderr(head, args...)
	return stdout, err
}

func (s LocalRunner) ExecGetStdoutStderr(head string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	//if !s.UnSudo {
	//	a := append([]string{}, head)
	//	args = append(a, args...)
	//	head = "sudo"
	//}

	if logs.IsDebugEnabled() {
		logs.WithField("command", strings.Join([]string{head, " ", strings.Join(args, " ")}, " ")).Debug("Running command")
	}
	cmd := exec.Command(head, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()
	err := cmd.Wait()
	return strings.TrimSpace(stdout.String()), stderr.String(), err
}
