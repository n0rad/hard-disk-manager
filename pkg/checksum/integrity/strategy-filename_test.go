package integrity

import (
	"github.com/n0rad/go-checksum/pkg/checksum"
)

var c = checksum.NewHash(checksum.Crc32_ieee)

//func TestTargetFile(t *testing.T) {
//	assert.Equal(t, "yopla/ouda-XX.txt", TargetFileReplacingSum(c, "yopla/ouda-2e0231fa.txt", "XX"))
//	assert.Equal(t, "yopla/ouda.xx-XX.txt", TargetFileReplacingSum(c, "yopla/ouda.xx-2e0231fa.txt", "XX"))
//	assert.Equal(t, "yopla/ouda-XX.txt", TargetFileReplacingSum(c, "yopla/ouda.txt", "XX"))
//	assert.Equal(t, "yopla/ouda-XX", TargetFileReplacingSum(c, "yopla/ouda", "XX"))
//	assert.Equal(t, "", TargetFileReplacingSum(c, "", "XX"))
//	assert.Equal(t, "a-XX", TargetFileReplacingSum(c, "a", "XX"))
//}
//
//func TestCRC32FromFilename(t *testing.T) {
//	assert.Equal(t, "2e0231fa", SumFromFilename(c, "yopla/ouda.genre-2e0231fa.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "-2e0231fa.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "-231fa.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "-245343231fa.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "yopla/ouda-2eX231fa.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "yopla/ouda.txt"))
//	assert.Equal(t, "", SumFromFilename(c, "a"))
//	assert.Equal(t, "", SumFromFilename(c, ""))
//}

//func TestCheckOrAddIntegrityFilenameSum(t *testing.T) {
//	sum, err := CheckOrAddIntegrityFilenameSum(c, "./README-bb3f3891.md")
//	assert.True(t, sum)
//	assert.Nil(t, err)
//}

//func TestCheckOrAddIntegrityFilenameSum(t *testing.T) {
//	Regex := regexp.MustCompile("(?i)\\.(socket|lock)$")
//	assert.NoError(t, directoryWalk("/tmp", Regex))
//}
