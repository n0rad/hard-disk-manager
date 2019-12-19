package hdm

import (
	"github.com/n0rad/hard-disk-manager/pkg/system"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisksSelector_Match(t *testing.T) {
	assert.True(t, DisksSelector{Server: "srv1"}.MatchDisk(Server{Name: "srv1"}, system.BlockDevice{}))
	//assert.False(t, DisksSelector{Server: "srv1"}.MatchDisk(Server{Name: "srv2"}, system.BlockDevice{BlockDeviceOLD: &system.BlockDevice{}}))
	//
	//assert.True(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Server{Name: "srv1"}, system.BlockDevice{BlockDeviceOLD: &system.BlockDevice{Partlabel: "toto"}}))
	//assert.False(t, DisksSelector{Server: "srv1", Label: "toto"}.MatchDisk(Server{Name: "srv1"}, system.BlockDevice{BlockDeviceOLD: &system.BlockDevice{Partlabel: "titi"}}))
	//
	//assert.True(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Server{Name: "srv1"}, system.BlockDevice{BlockDeviceOLD: &system.BlockDevice{Name: "toto"}}))
	//assert.False(t, DisksSelector{Server: "srv1", Disk: "toto"}.MatchDisk(Server{Name: "srv1"}, system.BlockDevice{BlockDeviceOLD: &system.BlockDevice{Name: "titi"}}))

}
