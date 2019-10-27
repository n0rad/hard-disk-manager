package runner

type Runner interface {
	ExecGetStdoutStderr(head string, args ...string) (string, string, error)
	ExecGetStdout(head string, args ...string) (string, error)
	ExecGetStd(head string, args ...string) (string, error)

	ExecShellGetStdout(cmd string) (string, error)
	ExecShellGetStd(cmd string) (string, error)
}
