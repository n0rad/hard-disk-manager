package system

import "strings"

// TODO get mounts devices
// systemctl show -p Where,Id *.mount --no-pager


// systemctl show -p What,Where,Id '*.mount' --no-pager | awk '/27b5d5103b8e.mount/' RS=


// find in mount units where this block device should be mounted
func SystemdMountPath(what string) (string, error) {
	s := Server{}
	s.Init()
	what = strings.Replace(what, `/`, `\/`, -1)
	return s.ExecShell("systemctl show -a -p What,Where,Id '*.mount' --no-pager | awk '/What=" + what + "/' RS= | grep Where | cut -f2 -d=")
}