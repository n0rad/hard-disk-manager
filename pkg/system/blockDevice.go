package system

import (
	"github.com/awnumar/memguard"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"strconv"
	"strings"
)

type BlockDevice struct {
	Fsavail    string        `json:"fsavail"`
	Fssize     string        `json:"fssize"`
	Fstype     string        `json:"fstype"`
	Fsused     string        `json:"fsused"`
	Fsuse      string        `json:"fsuse%"`
	Mountpoint string        `json:"mountpoint"`
	Label      string        `json:"label"`
	UUID       string        `json:"uuid"`
	Ptuuid     string        `json:"ptuuid"`
	Pttype     string        `json:"pttype"`
	Parttype   string        `json:"parttype"`
	Partlabel  string        `json:"partlabel"`
	Partuuid   string        `json:"partuuid"`
	Partflags  string        `json:"partflags"`
	Model      string        `json:"model"`
	Serial     string        `json:"serial"`
	State      string        `json:"state"`
	Owner      string        `json:"owner"`
	Group      string        `json:"group"`
	Wwn        string        `json:"wwn"`
	Pkname     string        `json:"pkname"`
	Hctl       string        `json:"hctl"`
	Tran       string        `json:"tran"`
	Rev        string        `json:"rev"`
	Vendor     string        `json:"vendor"`
	Name       string        `json:"name"`
	Kname      string        `json:"kname"`
	Path       string        `json:"path"`
	MajMin     string        `json:"maj:min"`
	Ra         int           `json:"ra"`
	Ro         bool          `json:"ro"`
	Rm         bool          `json:"rm"`
	Hotplug    bool          `json:"hotplug"`
	Size       string        `json:"size"`
	Mode       string        `json:"mode"`
	Alignment  int           `json:"alignment"`
	MinIo      int           `json:"min-io"`
	OptIo      int           `json:"opt-io"`
	PhySec     int           `json:"phy-sec"`
	LogSec     int           `json:"log-sec"`
	Rota       bool          `json:"rota"`
	Sched      string        `json:"sched"`
	RqSize     int           `json:"rq-size"`
	Type       string        `json:"type"`
	DiscAln    int           `json:"disc-aln"`
	DiscGran   string        `json:"disc-gran"`
	DiscMax    string        `json:"disc-max"`
	DiscZero   bool          `json:"disc-zero"`
	Wsame      string        `json:"wsame"`
	Rand       bool          `json:"rand"`
	Subsystems string        `json:"subsystems"`
	Zoned      string        `json:"zoned"`
	Children   []BlockDevice `json:"children"`

	fields data.Fields
	server *Server
}

// server can be nil
func (b *BlockDevice) Init(server *Server) {
	if server == nil {
		server = &Server{
		}
		_ = server.Init()
	}
	b.server = server
	b.fields = data.WithField("path", b.Path).WithField("server", b.server.Name)
	for i := range b.Children {
		b.Children[i].Init(server)
	}
}

func (b *BlockDevice) LuksOpen(cryptPassword *memguard.LockedBuffer) error {
	logs.WithFields(b.fields).Info("Disk luksOpen")
	if b.Fstype != "crypto_LUKS" {
		return errs.WithF(b.fields, "Cannot luks open, not a crypto block device")
	}

	if len(b.Children) > 0 {
		logs.WithFields(b.fields).Debug("Already open")
		return nil
	}

	volumeName := b.Partlabel
	if volumeName == "" {
		volumeName = b.Name
	}

	if out, err := b.server.ExecShell("echo -n '" + cryptPassword.String() + "' | sudo cryptsetup luksOpen " + b.Path + " " + volumeName + " -"); err != nil {
		return errs.WithEF(err, b.fields.WithField("out", out), "Failed to open luks")
	}

	return nil
}

func (b *BlockDevice) Mount(mountPath string) error {
	logs.WithFields(b.fields.WithField("mountPath", mountPath)).Debug("Mount")
	if mountPath == "" {
		return errs.WithF(b.fields, "mountPath cannot be empty")
	}

	if _, err := b.server.ExecShell("cat /proc/mounts | cut -f1,2 -d' ' | grep '" + b.Path + " " + mountPath + "$'"); err == nil {
		logs.WithF(b.fields).Debug("Directory is already mounted")
		return nil
	}

	if out, err := b.server.Exec("mkdir", "-p", mountPath); err != nil {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", string(out)), "Failed to create mount directory")
	}

	out, err := b.server.Exec("ls",  "-A", mountPath)
	if err != nil {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", string(out)), "Failed to ls on mount path")
	}
	if string(out) != "" {
		return errs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", string(out)), "Directory is not empty")
	}

	if out, err := b.server.ExecShell("! cat /proc/mounts | cut -f2 -d' ' | grep " + mountPath + "$"); err != nil {
		logs.WithEF(err, b.fields.WithField("path", mountPath).WithField("out", string(out))).Trace("Already mounted")
		return nil
	}

	if out, err := b.server.Exec("mount", b.Path, mountPath); err != nil {
		return errs.WithEF(err, b.fields.WithField("out", string(out)).WithField("target", mountPath), "Failed to mount")
	}
	return nil
}

func (b *BlockDevice) Umount(mountPath string) error {
	logs.WithFields(b.fields).Debug("Umount")

	if b.Mountpoint != "" {
		if out, err := b.server.Exec("umount", b.Mountpoint); err != nil {
			return errs.WithEF(err, b.fields.WithField("out", out), "Failed to unmount")
		}
	}

	if out, err := b.server.Exec("rmdir", mountPath); err != nil {
		logs.WithEF(err, b.fields.WithField("out", out)).Warn("Failed to cleanup mount path")
	}

	return nil
}

func (b *BlockDevice) SpaceAvailable() (int, error) {
	output, err := b.server.ExecShell("df " + b.Path + " --output=avail | tail -n +2")
	if err != nil {
		return 0, errs.WithEF(err, b.fields, "Failed to run 'df' on blockDevice")
	}

	size, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, errs.WithEF(err, b.fields.WithField("output", string(output)), "Failed to parse 'df' result")
	}
	return size, nil
}
