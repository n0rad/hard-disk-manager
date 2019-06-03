package hdm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisksSelector_Match(t *testing.T) {
	assert.True(t, DisksSelector{Server: "srv1"}.Match(Disk{ServerName: "srv1"}))
	assert.False(t, DisksSelector{Server: "srv1"}.Match(Disk{ServerName: "srv2"}))

	assert.True(t, DisksSelector{Server: "srv1", Label: "toto"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Partlabel: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Label: "toto"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Partlabel: "titi"}}))

	assert.True(t, DisksSelector{Server: "srv1", Disk: "toto"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Name: "toto"}}))
	assert.False(t, DisksSelector{Server: "srv1", Disk: "toto"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Name: "titi"}}))
	//assert.True(t, DisksSelector{Server: "srv1"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Partlabel: "toto"}}))
	//assert.True(t, DisksSelector{Server: "srv1"}.Match(Disk{ServerName: "srv1", BlockDevice: &BlockDevice{Name: "srv1"}}))


	//assert.True(t, DisksSelector{Server: "srv1", Label: "", Disk: "genre"}.Match(Disk{ServerName: "", BlockDevice: &BlockDevice{Partlabel: "", Name: "srv1"}}))
}
