package system

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/runner"
	"strings"
	"time"
)

const luksPartitionCode = "ca7d7ccb-63ed-4c53-861c-1742536059cc"

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
	exec   runner.Exec
}

func (b BlockDevice) String() string {
	return b.Path
}

func (b BlockDevice) HasChildren() bool {
	return len(b.Children) != 0
}

//func (b BlockDevice) GetFields() data.Fields {
//	return b.fields
//}

func (b BlockDevice) GetExec() runner.Exec {
	return b.exec
}

func (b *BlockDevice) Init(exec runner.Exec) {
	b.exec = exec
	b.fields = data.WithField("path", b.Path).WithField("exec", b.exec)

	for i := range b.Children {
		b.Children[i].Init(exec)
	}
}

func (b BlockDevice) LocationPath() (string, error) {
	output, err := b.exec.ExecGetStdout("find", "-L", "/dev/disk/by-path/", "-samefile", b.Path, "-printf", "%f")
	if err != nil {
		return "", errs.WithEF(err, b.fields, "Failed to get disk location")
	}
	return strings.TrimSpace(string(output)), nil
}

func (b *BlockDevice) Reload() error {
	time.Sleep(1000 * time.Millisecond) // TODO info are mising when lsblk is run just after change

	lsblk := Lsblk{}
	if err := lsblk.Init(b.exec); err != nil {
		return errs.WithE(err, "Failed to init lsblk to reload blockDevice")
	}

	device, err := lsblk.GetBlockDevice(b.Path)
	if err != nil {
		return errs.WithE(err, "Failed to get device fro; lsblk to reload")
	}

	*b = device
	return nil
}
