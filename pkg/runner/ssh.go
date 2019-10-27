package runner

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/sfreiberg/simplessh"
	"strings"
)

type SshExec struct {
	Hostname string
	Username string

	sshClient *simplessh.Client
}

func (s SshExec) ExecGetStdoutStderr(head string, args ...string) (string, string, error) {
	if s.sshClient == nil {
		client, err := simplessh.ConnectWithAgent(s.Hostname, s.Username)
		if err != nil {
			return "", "", errs.WithEF(err, data.WithField("hostname", s.Hostname).WithField("username", s.Username), "Fail to ssh to server")
		}
		s.sshClient = client
	}

	cmd := strings.Join([]string{head, " ", strings.Join(args, " ")}, " ")
	logs.WithField("host", s.Hostname).WithField("cmd", cmd).Debug("Running command on server")

	stdout, err := s.sshClient.Exec(cmd)
	logs.WithField("stdout", string(stdout)).WithField("command", cmd).Trace("command output")
	if err != nil {
		return string(stdout), "", errs.WithEF(err, data.WithField("host", s.Hostname).
			WithField("cmd", cmd), "Exec command failed")
	}

	return strings.TrimSpace(string(stdout)), "", nil
}

func (s SshExec) ExecGetStdout(head string, args ...string) (string, error) {
	stdout, _, err := s.ExecGetStdoutStderr(head, args...)
	return stdout, err
}

func (s SshExec) ExecGetStd(head string, args ...string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr(head, args...)
	stdout += stderr
	return stdout, err
}

/////////////////

func (s SshExec) ExecShellGetStdout(cmd string) (string, error) {
	stdout, stderr, err := s.ExecGetStdoutStderr("bash", "-o", "pipefail", "-c", cmd)
	stdout += stderr
	return stdout, err
}

func (s SshExec) ExecShellGetStd(cmd string) (string, error) {
	return s.ExecGetStdout("bash", "-o", "pipefail", "-c", cmd)
}

func (s *SshExec) Close() {
	if s.sshClient != nil {
		s.sshClient.Close()
		s.sshClient = nil
	}
}
