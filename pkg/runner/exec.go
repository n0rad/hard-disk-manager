package runner

import "io"

type Exec interface {
	ExecSetStdinGetStdoutStderr(stdin io.Reader, head string, args ...string) (string, string, error)
	ExecGetStdoutStderr(head string, args ...string) (string, string, error)
	ExecGetStdout(head string, args ...string) (string, error)
	ExecGetStd(head string, args ...string) (string, error)
	Exec(head string, args ...string) error

	ExecShellGetStdout(cmd string) (string, error)
	ExecShellGetStd(cmd string) (string, error)

	Close()
}
