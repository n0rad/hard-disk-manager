package system

type Lsblk struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
}
