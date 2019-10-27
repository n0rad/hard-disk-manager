package system

import (
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"strings"
)

// TODO get mounts devices
// systemctl show -p Where,Id *.mount --no-pager

// systemctl show -p What,Where,Id '*.mount' --no-pager | awk '/27b5d5103b8e.mount/' RS=

// find in mount units where this block device should be mounted
type Systemd struct {
	exec runner.Exec
}

func (s *Systemd) Init(exec runner.Exec) {
	s.exec = exec
}

func (s *Systemd) SystemdMountPath(what string) (string, error) {
	what = strings.Replace(what, `/`, `\/`, -1)
	return s.exec.ExecShellGetStdout("systemctl show -a -p What,Where,Id '*.mount' --no-pager | awk '/What=" + what + "/' RS= | grep Where | cut -f2 -d=")
}
