package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisksSelector_Match(t *testing.T) {
	assert.True(t, DisksSelector{Server: "srv1"}.MatchDisk(Disk{ServerName: "srv1", BlockDeviceOLD: &BlockDeviceOLD{}}))
	assert.False(t, DisksSelector{Server: "srv1"}.MatchDisk(Disk{ServerName: "srv2", BlockDeviceOLD: &BlockDeviceOLD{}}))

	assert.True(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDeviceOLD: &BlockDeviceOLD{Partlabel: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDeviceOLD: &BlockDeviceOLD{Partlabel: "titi"}}))

	assert.True(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDeviceOLD: &BlockDeviceOLD{Name: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDeviceOLD: &BlockDeviceOLD{Name: "titi"}}))

}
