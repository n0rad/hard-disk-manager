package bench

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"
	"os"
	"testing"
)

//
//func benchmarkRun(h hash.Hash, i int, b *testing.B) (string, error) {
//	file, err := os.Open("/home/n0rad/Tmp/Anno.1800.Crack.Only.rar")
//	if err != nil {
//		b.Fatal(err)
//	}
//	defer file.Close()
//
//	bs := make([]byte, i)
//
//	for {
//		// read a chunk
//		n, err := file.Read(bs)
//		if err != nil && err != io.EOF {
//			panic(err)
//		}
//		if n == 0 {
//			break
//		}
//
//		// write a chunk
//		if _, err := h.Write(bs[:n]); err != nil {
//			panic(err)
//		}
//	}
//
//	return hex.EncodeToString(h.Sum(nil)), nil
//}
//
//func BenchmarkMD5_1k(b *testing.B) {
//	benchmarkRun(md5.New(), 1024, b)
//}
//
//func BenchmarkMD5_10k(b *testing.B) {
//	benchmarkRun(md5.New(), 10*1024, b)
//}
//
//func BenchmarkMD5_100k(b *testing.B) {
//	benchmarkRun(md5.New(), 100*1024, b)
//}
//
//func BenchmarkMD5_250k(b *testing.B) {
//	benchmarkRun(md5.New(), 250*1024, b)
//}
//
//func BenchmarkMD5_500k(b *testing.B) {
//	benchmarkRun(md5.New(), 500*1024, b)
//}
//
//func BenchmarkSHA11_1k(b *testing.B) {
//	benchmarkRun(sha1.New(), 1024, b)
//}
//
//func BenchmarkSha1_10k(b *testing.B) {
//	benchmarkRun(sha1.New(), 10*1024, b)
//}
//
//func BenchmarkSha1_100k(b *testing.B) {
//	benchmarkRun(sha1.New(), 100*1024, b)
//}
//
//func BenchmarkSha1_250k(b *testing.B) {
//	benchmarkRun(sha1.New(), 250*1024, b)
//}
//
//func BenchmarkSha1_500k(b *testing.B) {
//	benchmarkRun(sha1.New(), 500*1024, b)
//}
//
//func BenchmarkSha256_1k(b *testing.B) {
//	benchmarkRun(sha256.New(), 1024, b)
//}
//
//func BenchmarkSha256_10k(b *testing.B) {
//	benchmarkRun(sha256.New(), 10*1024, b)
//}
//
//func BenchmarkSha256_100k(b *testing.B) {
//	benchmarkRun(sha256.New(), 100*1024, b)
//}
//
//func BenchmarkSha256_250k(b *testing.B) {
//	benchmarkRun(sha256.New(), 250*1024, b)
//}
//
//func BenchmarkCRC32_1k(b *testing.B) {
//	benchmarkRun(sha256.New(), 1024, b)
//}
//
//func BenchmarkCRC32_10k(b *testing.B) {
//	benchmarkRun(sha256.New(), 10*1024, b)
//}
//
//func BenchmarkCRC32_100k(b *testing.B) {
//	benchmarkRun(sha256.New(), 100*1024, b)
//}
//
//func BenchmarkCRC32_250k(b *testing.B) {
//	benchmarkRun(sha256.New(), 250*1024, b)
//}
//
//func BenchmarkCRC32_500k(b *testing.B) {
//	benchmarkRun(sha256.New(), 500*1024, b)
//}
//
//func BenchmarkCRC32_xk(b *testing.B) {
//	sumFile(sha256.New(), "/home/n0rad/Tmp/Anno.1800.Crack.Only.rar")
//}

const file = "/home/n0rad/Tmp/test.mkv"

func BenchmarkFnv32(b *testing.B) {
	sumFile(fnv.New32(), file)
}

func BenchmarkFnv128(b *testing.B) {
	sumFile(fnv.New128(), file)
}

func BenchmarkMd5(b *testing.B) {
	sumFile(md5.New(), file)
}

func BenchmarkSha1(b *testing.B) {
	sumFile(sha1.New(), file)
}

func BenchmarkSha256(b *testing.B) {
	sumFile(sha256.New(), file)
}

func BenchmarkSha512(b *testing.B) {
	sumFile(sha512.New(), file)
}

func BenchmarkCRC32(b *testing.B) {
	sumFile(crc32.New(crc32.IEEETable), file)
}

func BenchmarkCRC64iso(b *testing.B) {
	sumFile(crc64.New(crc64.MakeTable(crc64.ISO)), file)
}

func BenchmarkCRC64ecma(b *testing.B) {
	sumFile(crc64.New(crc64.MakeTable(crc64.ECMA)), file)
}

func sumFile(hash hash.Hash, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

//func hash_file_crc32(filePath string, polynomial uint32) (string, error) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return "", err
//	}
//	defer file.Close()
//
//	tablePolynomial := crc32.MakeTable(polynomial)
//	hash := crc32.New(tablePolynomial)
//
//
//	//Copy the file in the interface
//	if _, err := io.Copy(hash, file); err != nil {
//		return "", err
//	}
//
//	//Generate the hash
//	hashInBytes := hash.Sum(nil)[:]
//
//	returnCRC32String := hex.EncodeToString(hashInBytes)
//
//	//Return the output
//	return returnCRC32String, nil
//}

//
//
//func makeHash(name string) hash.Hash {
//	switch strings.ToLower(name) {
//	case "blake2b-256":
//		return mustMakeHash(blake2b.New256(nil))
//	case "blake2b-384":
//		return mustMakeHash(blake2b.New384(nil))
//	case "blake2b-512":
//		return mustMakeHash(blake2b.New512(nil))
//	case "blake2s-256":
//		return mustMakeHash(blake2s.New256(nil))
//	case "ripemd160":
//		return ripemd160.New()
//	case "md4":
//		return md4.New()
//	case "md5":
//		return md5.New()
//	case "sha1":
//		return sha1.New()
//	case "sha256":
//		return sha256.New()
//	case "sha384":
//		return sha512.New384()
//	case "sha3-224":
//		return sha3.New224()
//	case "sha3-256":
//		return sha3.New256()
//	case "sha3-384":
//		return sha3.New384()
//	case "sha3-512":
//		return sha3.New512()
//	case "sha512":
//		return sha512.New()
//	case "sha512-224":
//		return sha512.New512_224()
//	case "sha512-256":
//		return sha512.New512_256()
//	case "crc32-ieee":
//		return crc32.NewIEEE()
//	case "crc64-iso":
//		return crc64.New(crc64.MakeTable(crc64.ISO))
//	case "crc64-ecma":
//		return crc64.New(crc64.MakeTable(crc64.ECMA))
//	case "adler32":
//		return adler32.New()
//	case "fnv32":
//		return fnv.New32()
//	case "fnv32a":
//		return fnv.New32a()
//	case "fnv64":
//		return fnv.New64()
//	case "fnv64a":
//		return fnv.New64a()
//	case "fnv128":
//		return fnv.New128()
//	case "fnv128a":
//		return fnv.New128a()
//	case "xor8":
//		return new(xor8)
//	case "fletch16":
//		return &fletch16{}
//	case "luhn":
//		return new(luhn)
//	case "sum16":
//		return new(sum16)
//	case "sum32":
//		return new(sum32)
//	case "sum64":
//		return new(sum64)
//	case "crc8":
//		return new(crc8)
//	case "crc16-ccitt":
//		c := new(crc16ccitt)
//		c.Reset()
//		return c
//	}
//	return nil
//}
