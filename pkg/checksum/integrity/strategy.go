package integrity

import "github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"

type Strategy interface {
	IsSet(file string) (bool, error)
	GetSum(file string) (string, error)
	Sum(file string) (string, error)       // TODO generic
	SumAndSet(file string) (string, error) // TODO generic
	Set(file string, sum string) error
	Remove(file string) error
	Check(file string) (error, error) // TODO generic
	IsSumFile(file string) bool
}

func NewStrategy(strategyName string, hash hashs.Hash) Strategy {
	switch strategyName {
	case "sumfile":
		return StrategySumFile{
			Hash:     hashs.NewHash(hash),
			HashName: string(hash),
		}
	case "filename":
		return StrategyFilename{
			Hash:    hashs.NewHash(hash),
			OldHash: hashs.NewHash(hash), // support old HASH
		}
	}
	return nil
}
