package hdm

type Disks []Disk

func (d Disks) findDiskBySelector(selector DisksSelector) *Disk {
	for _, disk := range d {
		if selector.MatchDisk(disk) {
			return &disk
		}
	}
	return nil
}

func (d Disks) findDeepestBlockDeviceByLabel(label string) *BlockDevice {
	if label == "" {
		return nil
	}
	for _, disk := range d {
		for _, partition := range disk.Children {
			if partition.Partlabel == label {
				device := findDeepestBlockDevice(partition)
				return &device
			}
		}
	}
	return nil
}

func findDeepestBlockDevice(device BlockDevice) BlockDevice {
	if len(device.Children) > 0 {
		return findDeepestBlockDevice(device.Children[0])
	}
	return device
}
