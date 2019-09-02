package system

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/tools"
	"strings"
)

//partprobe
//wipefs --all /dev/sdX
//sudo lsblk -o name,size,type,fstype,label,partlabel,mountpoint,path

type Disk struct {
	*BlockDevice
	SmartResult *tools.SmartResult

	ServerName string `json:"server"`
}

func (b *Disk) String() string {
	return b.Path
}

func (d *Disk) Init(server *Server) {
	d.BlockDevice.Init(server, d)
}

func (d *Disk) LocationPath() (string, error) {
	output, err := d.server.Exec("find -L /dev/disk/by-path/ -samefile " + d.Path + " -printf '%f\n'")
	if err != nil {
		return "", errs.WithEF(err, d.fields, "Failed to get disk location")
	}
	return strings.TrimSpace(string(output)), nil
}

func (d *Disk) Location() (string, error) {
	path, err := d.LocationPath()
	if err != nil {
		return "", err
	}
	return d.server.BayLocation(path), nil
}

func (d *Disk) Add(password string) error {
	for {
		newDevices, err := d.addAndGiveNewDevices(password)
		if err != nil {
			return err
		}
		if newDevices == false {
			break
		} else if err := d.ReplaceFromLsblk(); err != nil {
			return errs.WithEF(err, d.fields, "Failed to scan after luksOpen")
		}
	}
	return nil
}

func (d *Disk) HasChildren() bool {
	if len(d.Children) == 0 {
		return false
	}
	return true
}

func (d *Disk) PopulateFromBytes(bytes []byte) error {
	if err := yaml.Unmarshal(bytes, d); err != nil {
		return errs.WithEF(err, data.WithField("data", string(bytes)), "Failed to unmarshal disk db file")
	}

	// TODO setup runner & fields

	return nil
}

func (d *Disk) ReplaceFromLsblk() error {
	logs.WithFields(d.fields).Info("Running lsblk")

	output, err := d.server.Exec("sudo lsblk -J -O " + d.Path)
	if err != nil {
		return errs.WithE(err, "Fail to get disk from lsblk")
	}

	lsblk := Lsblk{}
	if err = json.Unmarshal([]byte(output), &lsblk); err != nil {
		return errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal lsblk result")
	}

	lsblk.Blockdevices[0].Init(d.server)
	*d = lsblk.Blockdevices[0]
	return nil
}

func (d *Disk) FillFromSmartctl() error {
	logs.WithFields(d.fields).Info("Running smartctl")

	output, err := d.server.Exec("sudo smartctl --xall -j " + d.Path + " || true")
	if err != nil {
		return errs.WithEF(err, d.fields, "Fail to run smartctl")
	}
	logs.WithField("output", string(output)).Trace("smart output")

	smartResult := tools.SmartResult{}
	if err = json.Unmarshal([]byte(output), &smartResult); err != nil {
		return errs.WithEF(err, data.WithField("payload", string(output)), "Fail to unmarshal smartctl result")
	}
	d.SmartResult = &smartResult

	//	if attribute.ID == 5 {
	//		d.ReallocatedSectorCount = attribute.Raw.Value
	//	if attribute.ID == 187 {
	//		d.ReportedUncorrect = attribute.Raw.Value
	//	if attribute.ID == 188 {
	//		d.CommandTimeout = attribute.Raw.Value
	//	if attribute.ID == 197 {
	//		d.CurrentPendingSector = attribute.Raw.Value
	//	if attribute.ID == 198 {
	//		d.OfflineUncorrectable = attribute.Raw.Value
	return nil
}

func (d *Disk) Scan() error {
	if err := d.ReplaceFromLsblk(); err != nil {
		return errs.WithEF(err, d.fields, "Fail to rescan disk after luksFormat")
	}

	if err := d.FillFromSmartctl(); err != nil {
		return errs.WithEF(err, d.fields, "Failed to add smartctl info disk")
	}
	return nil
}

func (d *Disk) Prepare(label string, cryptPassword string) error {
	if len(d.Children) != 0 {
		return errs.WithF(d.fields, "Cannot prepare disk, some partitions exists")
	}

	logs.WithFields(d.fields.WithField("label", label)).Info("Prepare disk")

	_, err := d.server.Exec("sudo sgdisk -og " + d.Path)
	if err != nil {
		return errs.WithEF(err, d.fields, "Fail to clear partition table")
	}

	_, err = d.server.Exec(`sudo sgdisk -n 1:0:0 -t 1:CA7D7CCB-63ED-4C53-861C-1742536059CC -c 1:"` + label + `" ` + d.Path)
	if err != nil {
		return errs.WithEF(err, d.fields, "Fail to create partition")
	}

	if err := d.Scan(); err != nil {
		return errs.WithEF(err, d.fields, "Fail to rescan disk after luksFormat")
	}

	if len(d.Children) != 1 {
		return errs.WithF(d.fields, "Number of partitions is not one after prepare")
	}

	if _, err = d.server.Exec("echo -n '" + cryptPassword + "' | sudo cryptsetup --verbose --hash=sha512 --cipher=aes-xts-benbi:sha512 --key-size=512 luksFormat " + d.Children[0].Path + " -"); err != nil {
		return errs.WithEF(err, d.fields, "Fail to crypt partition")
	}

	if err := d.Scan(); err != nil {
		return errs.WithEF(err, d.fields, "Failed to rescan disk after luksFormat")
	}

	if err := d.Children[0].luksOpen(cryptPassword); err != nil {
		return errs.WithEF(err, d.fields, "Failed to open crypt partition")
	}

	if err := d.Scan(); err != nil {
		return errs.WithEF(err, d.fields, "Failed to rescan disk after luksOpen")
	}

	if _, err = d.server.Exec("sudo mkfs.xfs -L " + label + " -f " + d.Children[0].Children[0].Path); err != nil {
		return errs.WithEF(err, d.fields, "Failed to make filesystem")
	}

	if err := d.Scan(); err != nil {
		return errs.WithEF(err, d.fields, "Failed to rescan disk after luksOpen")
	}

	if err := d.Children[0].Children[0].luksClose(); err != nil {
		return errs.WithEF(err, d.fields, "Failed to close partition")
	}

	return nil
}

func (d *Disk) Erase() error {
	if len(d.Children) > 0 {
		return errs.WithF(d.fields, "Disk has partitions")
	}
	if _, err := d.server.Exec("sudo hdparm --user-master u --security-erase Nine " + d.Path); err != nil {
		return errs.WithEF(err, d.fields, "Fail to erase disk")
	}
	return nil
}

// https://www.vincentliefooghe.net/content/linux-corriger-des-secteurs-d%C3%A9fecteux-sur-un-disque
func (d *Disk) smartRepairPendingSectors() {

}
