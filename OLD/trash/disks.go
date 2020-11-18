package trash

type Disks []Disk

func (d Disks) findDiskBySelector(selector DisksSelector) *Disk {
	for _, disk := range d {
		if selector.MatchDisk(disk) {
			return &disk
		}
	}
	return nil
}

func (d Disks) FindDeepestBlockDeviceByLabel(label string) *BlockDeviceOLD {
	if label == "" {
		return nil
	}
	for _, disk := range d {
		for _, partition := range disk.Children {
			if partition.Partlabel == label {
				device := partition.FindDeepestBlockDevice()
				return &device
			}
		}
	}
	return nil
}
