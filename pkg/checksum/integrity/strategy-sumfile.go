package integrity

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
)

type StrategySumFile struct {
	Hash     hash.Hash
	HashName string
}

func (s StrategySumFile) IsSumFile(file string) bool {
	if filepath.Ext(file) == s.sumFileSuffix() {
		return true
	}
	return false
}

func (s StrategySumFile) IsSet(file string) (bool, error) {
	if _, err := os.Stat(s.sumFilename(file)); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (s StrategySumFile) GetSum(file string) (string, error) {
	sumFilename := s.sumFilename(file)
	readFile, err := ioutil.ReadFile(sumFilename)
	if err != nil {
		return "", errs.WithEF(err, data.WithField("file", sumFilename), "Failed to read sum file")
	}
	if s.HashName == string(hashs.Crc32_ieee) {
		return hashs.SumFromSumSfvLine(string(readFile)), nil
	} else {
		return hashs.SumFromSumLine(string(readFile)), nil
	}
}

func (s StrategySumFile) Sum(file string) (string, error) {
	sum, err := hashs.SumFile(s.Hash, file)
	if err != nil {
		return "", errs.WithEF(err, data.WithField("file", file), "Failed to sum file")
	}
	return sum, nil
}

func (s StrategySumFile) SumAndSet(file string) (string, error) {
	sum, err := hashs.SumFile(s.Hash, file)
	if err != nil {
		return "", err
	}
	return sum, s.Set(file, sum)
}

func (s StrategySumFile) Set(file string, sum string) error {
	stat, err := os.Stat(file)
	if err != nil {
		return errs.WithEF(err, data.WithField("file", file), "Failed to get stat of file")
	}

	var line string
	if s.HashName == string(hashs.Crc32_ieee) {
		line = hashs.SumSfvLine(filepath.Base(file), sum)
	} else {
		line = hashs.SumLine(filepath.Base(file), sum)
	}

	sumFilename := s.sumFilename(file)
	if err := ioutil.WriteFile(sumFilename, []byte(line), stat.Mode()); err != nil {
		return errs.WithEF(err, data.WithField("file", sumFilename), "Failed to write sum file")
	}
	return nil
}

func (s StrategySumFile) Remove(file string) error {
	sumFilename := s.sumFilename(file)
	if _, err := os.Stat(sumFilename); err == nil {
		if err := os.Remove(sumFilename); err != nil {
			return errs.WithEF(err, data.WithField("file", sumFilename), "Failed to remove sum file")
		}
	} else if !os.IsNotExist(err) {
		return errs.WithF(data.WithField("file", sumFilename), "Failed to stat sum file")
	}
	return nil
}

func (s StrategySumFile) Check(file string) (error, error) {
	sum, err := s.Sum(file)
	if err != nil {
		return nil, err
	}

	savedSum, err := s.GetSum(file)
	if err != nil {
		return nil, errs.WithE(err, "Failed to get saved sum")
	}
	if savedSum != sum {
		return errs.WithF(data.WithField("sum", sum).WithField("saved-sum", savedSum), "Sums do not match"), nil
	}
	return nil, nil
}

////////////////////////////////////////

func (s StrategySumFile) sumFilename(file string) string {
	return file + s.sumFileSuffix()
}

func (s StrategySumFile) sumFileSuffix() string {
	if s.HashName == string(hashs.Crc32_ieee) {
		return ".sfv"
	}
	return "." + s.HashName
}
