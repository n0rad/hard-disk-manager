package hashs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/algo"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
)

func MakeHashString(hashName string) hash.Hash {
	return NewHash(Hash(hashName))
}

func NewHash(hash Hash) hash.Hash {
	switch hash {
	case "blake2b-256":
		new256, _ := blake2b.New256(nil)
		return new256
	case "blake2b-384":
		new384, _ := blake2b.New384(nil)
		return new384
	case "blake2b-512":
		new512, _ := blake2b.New512(nil)
		return new512
	case "blake2s-256":
		new256, _ := blake2s.New256(nil)
		return new256
	case "ripemd160":
		return ripemd160.New()
	case "md4":
		return md4.New()
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha384":
		return sha512.New384()
	case "sha3-224":
		return sha3.New224()
	case "sha3-256":
		return sha3.New256()
	case "sha3-384":
		return sha3.New384()
	case "sha3-512":
		return sha3.New512()
	case "sha512":
		return sha512.New()
	case "sha512-224":
		return sha512.New512_224()
	case "sha512-256":
		return sha512.New512_256()
	case "crc32":
		return crc32.NewIEEE()
	case "crc64-iso":
		return crc64.New(crc64.MakeTable(crc64.ISO))
	case "crc64-ecma":
		return crc64.New(crc64.MakeTable(crc64.ECMA))
	case "adler32":
		return adler32.New()
	case "fnv32":
		return fnv.New32()
	case "fnv32a":
		return fnv.New32a()
	case "fnv64":
		return fnv.New64()
	case "fnv64a":
		return fnv.New64a()
	case "fnv128":
		return fnv.New128()
	case "fnv128a":
		return fnv.New128a()
	case "xor8":
		return new(algo.Xor8)
	case "fletch16":
		return &algo.Fletch16{}
	case "luhn":
		return new(algo.Luhn)
	case "sum16":
		return new(algo.Sum16)
	case "sum32":
		return new(algo.Sum32)
	case "sum64":
		return new(algo.Sum64)
	case "crc8":
		return new(algo.Crc8)
	case "crc16-ccitt":
		c := new(algo.Crc16ccitt)
		c.Reset()
		return c
	}
	return nil
}
