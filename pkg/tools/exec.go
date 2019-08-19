package tools

import (
	"bytes"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/sfreiberg/simplessh"
	"os/exec"
	"strings"
)

type Runner interface {
	ExecGetOutputError(head string, args ...string) (string, string, error)
}

type SshRunner struct {
	Hostname string
	Username string

	sshClient *simplessh.Client
}

func (s *SshRunner) Close() {
	if s.sshClient != nil {
		s.sshClient.Close()
		s.sshClient = nil
	}
}

func (s *SshRunner) ExecGetOutputError(head string, args ...string) (string, string, error) {
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

//

type LocalRunner struct {
	UnSudo bool
}

func (s LocalRunner) ExecGetOutputError(head string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if !s.UnSudo {
		a := append([]string{}, head)
		args = append(a, args...)
		head = "sudo"
	}

	if logs.IsDebugEnabled() {
		logs.WithField("command", strings.Join([]string{head, " ", strings.Join(args, " ")}, " ")).Debug("Running command")
	}
	cmd := exec.Command(head, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()
	err := cmd.Wait()
	return strings.TrimSpace(stdout.String()),  stderr.String(), err
}
