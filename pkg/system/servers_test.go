package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisksSelector_Match(t *testing.T) {
	assert.True(t, DisksSelector{Server: "srv1"}.MatchDisk(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{}}))
	assert.False(t, DisksSelector{Server: "srv1"}.MatchDisk(Disk{ServerName: "srv2", BlockDevice: &BlockDevice{}}))

	assert.True(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Partlabel: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Partlabel: "titi"}}))

	assert.True(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Name: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Name: "titi"}}))

}
