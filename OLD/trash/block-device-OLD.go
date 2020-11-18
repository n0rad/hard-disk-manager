package trash

import (
	"github.com/awnumar/memguard"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/utils"
)

type LsblkOLD struct {
	Blockdevices []Disk `json:"blockdevices"`
}

type BlockDeviceOLD struct {
	Fsavail    string           `json:"fsavail"`
	Fssize     string           `json:"fssize"`
	Fstype     string           `json:"fstype"`
	Fsused     string           `json:"fsused"`
	Fsuse      string           `json:"fsuse%"`
	Mountpoint string           `json:"mountpoint"`
	Label      string           `json:"label"`
	UUID       string           `json:"uuid"`
	Ptuuid     string           `json:"ptuuid"`
	Pttype     string           `json:"pttype"`
	Parttype   string           `json:"parttype"`
	Partlabel  string           `json:"partlabel"`
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
	Rota       bool             `json:"rota"`
	Sched      string           `json:"sched"`
	RqSize     int              `json:"rq-size"`
	Type       string           `json:"type"`
	DiscAln    int              `json:"disc-aln"`
	DiscGran   string           `json:"disc-gran"`
	DiscMax    string           `json:"disc-max"`
	DiscZero   bool             `json:"disc-zero"`
	Wsame      string           `json:"wsame"`
	Rand       bool             `json:"rand"`
	Subsystems string           `json:"subsystems"`
	Zoned      string           `json:"zoned"`
	Children   []BlockDeviceOLD `json:"children"`

	//disk *Disk
	server *Server
	fields data.Fields
}

func (b *BlockDeviceOLD) String() string {
	return b.Path
}

func (b *BlockDeviceOLD) Init(server *Server, disk *Disk) {
	b.server = server
	//b.disk = disk
	b.fields = data.WithField("path", b.Path).WithField("server", b.server.Name)
	for i := range b.Children {
		b.Children[i].Init(server, disk)
	}
}

func (b BlockDeviceOLD) FindDeepestBlockDevice() BlockDeviceOLD {
	if len(b.Children) > 0 {
		return b.Children[0].FindDeepestBlockDevice()
	}
	return b
}

func (b *BlockDeviceOLD) addAndGiveNewDevices(password *memguard.LockedBuffer) (bool, error) {
	logs.WithFields(b.fields).Debug("Disk add")
	if len(b.Children) > 0 {
		newDevices := false
		for _, child := range b.Children {
			newRecursive, err := child.addAndGiveNewDevices(password)
			if err != nil {
				logs.WithEF(err, b.fields).Warn("Cannot add device")
			}
			if newRecursive == true {
				newDevices = newRecursive
			}
		}
		return newDevices, nil
	}

	newDevices := false
	if b.Fstype == "crypto_LUKS" {
		if err := b.luksOpen(password); err != nil {
			return false, err
		}
		newDevices = true
	} else if utils.SliceContains(filesystems, b.Fstype) {
		if err := b.mount(); err != nil {
			b.deleteMountDir()
			return false, err
		}
	} else {
		return false, errs.WithF(b.fields.WithField("fstype", b.Fstype), "Unknown fstype")
	}
	return newDevices, nil
}

func (b *BlockDeviceOLD) Remove() error {
	logs.WithFields(b.fields).Info("Disk remove")
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			if err := child.Remove(); err != nil {
				logs.WithE(err).Warn("Cannot remove device")
			}
		}
	}

	if b.Mountpoint != "" {
		if err := b.unmount(); err != nil {
			return err
		}
	}

	if utils.SliceContains(filesystems, b.Fstype) {
		b.deleteMountDir()
	}

	switch b.Type {
	case "crypt":
		if err := b.luksClose(); err != nil {
			return err
		}
	case "disk":
		if err := b.putInSleepNow(); err != nil {
			return err
		}
	}
	return nil
}
