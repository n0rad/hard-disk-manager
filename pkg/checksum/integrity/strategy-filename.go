package integrity

import (
	"encoding/hex"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"
	"hash"
	"os"
	"path/filepath"
	"strings"
)

type StrategyFilename struct {
	Hash    hash.Hash
	OldHash hash.Hash
}

func (s StrategyFilename) IsSumFile(file string) bool {
	return false
}

func (s StrategyFilename) Set(file string, sum string) error {
	filename := s.newFilename(file, sum)
	if err := os.Rename(file, filename); err != nil {
		return errs.WithEF(err, data.WithField("new-file", filename), "Failed to rename file to set sum")
	}
	return nil
}

func (s StrategyFilename) SumAndSet(file string) (string, error) {
	sum, err := hashs.SumFile(s.Hash, file)
	if err != nil {
		return "", err
	}
	return sum, s.Set(file, sum)
}

func (s StrategyFilename) GetSum(file string) (string, error) {
	fileWithoutExt := strings.TrimSuffix(file, filepath.Ext(file))
	hHexLen := s.Hash.Size() * 2
	if len(filepath.Base(fileWithoutExt))-hHexLen <= 1 { // filename only contains CRC ?
		return "", nil
	}
	if (fileWithoutExt[len(fileWithoutExt)-hHexLen-1]) != '-' { // crc do not start with a -
		return "", nil
	}
	candidate := fileWithoutExt[len(fileWithoutExt)-hHexLen:]
	_, err := hex.DecodeString(candidate) // not a crc
	if err != nil {
		return "", nil
	}
	return candidate, nil
}

func (s StrategyFilename) IsSet(file string) (bool, error) {
	sum, err := s.GetSum(file)
	if err != nil {
		return false, err
	}
	if sum == "" {
		return false, nil
	}
	return true, nil
}

func (s StrategyFilename) Sum(file string) (string, error) {
	sum, err := hashs.SumFile(s.Hash, file)
	if err != nil {
		return "", errs.WithEF(err, data.WithField("file", file), "Failed to sum file")
	}
	return sum, nil
}

func (s StrategyFilename) Remove(file string) error {
	return s.Set(file, "")
}

func (s StrategyFilename) Check(file string) (error, error) {
	sum, err := s.Sum(file)
	if err != nil {
		return nil, err
	}

	savedSum, err := s.GetSum(file)
	if err != nil {
		return nil, errs.WithE(err, "Failed to get saved sum")
	}
	if savedSum != sum {
		return errs.WithF(data.WithField("sum", sum).WithField("saved-sum", savedSum), "sums do not match"), nil
	}
	return nil, nil
}

//////////////////////////////

func (s StrategyFilename) newFilename(file string, newSum string) string {
	if file == "" {
		return ""
	}

	b := strings.Builder{}
	pathBase := strings.TrimSuffix(file, filepath.Ext(file))
	sum, _ := s.GetSum(file) // cannot return an error
	if sum != "" {
		b.WriteString(pathBase[:len(pathBase)-(s.OldHash.Size()*2)-1])
	} else {
		b.WriteString(pathBase)
	}
	if newSum != "" {
		b.WriteRune('-')
		b.WriteString(newSum)
	}
	b.WriteString(filepath.Ext(file))
	return b.String()
}
