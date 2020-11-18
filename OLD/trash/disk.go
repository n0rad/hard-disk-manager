package trash

import (
	"encoding/json"
	"github.com/awnumar/memguard"
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"strings"
)

//partprobe
//wipefs --all /dev/sdX
//sudo lsblk -o name,size,type,fstype,label,partlabel,mountpoint,path

type Disk struct {
	*BlockDeviceOLD
	SmartResult *SmartResult

	ServerName string `json:"server"`
}

func (b *Disk) String() string {
	return b.Path
}

func (d *Disk) Init(server *Server) {
	d.BlockDeviceOLD.Init(server, d)
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

func (d *Disk) Add(password *memguard.LockedBuffer) error {
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

func (d *Disk) PopulateFromBytes(bytes []byte) error {
	if err := yaml.Unmarshal(bytes, d); err != nil {
		return errs.WithEF(err, data.WithField("data", string(bytes)), "Failed to unmarshal disk db file")
	}

	// TODO setup runner & fields

	return nil
}

func (d *Disk) FillFromSmartctl() error {
	logs.WithFields(d.fields).Info("Running smartctl")

	output, err := d.server.Exec("sudo smartctl --xall -j " + d.Path + " || true")
	if err != nil {
		return errs.WithEF(err, d.fields, "Fail to run smartctl")
	}
	logs.WithField("output", string(output)).Trace("smart output")

	smartResult := SmartResult{}
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


// https://www.vincentliefooghe.net/content/linux-corriger-des-secteurs-d%C3%A9fecteux-sur-un-disque
func (d *Disk) smartRepairPendingSectors() {

}
