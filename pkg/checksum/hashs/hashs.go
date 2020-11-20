package hashs

type Hash string

const Blake2b_256 Hash = "blake2b-256"
const Blake2b_384 Hash = "blake2b-384"
const Blake2b_512 Hash = "blake2b-512"
const Blake2s_256 Hash = "blake2s-256"
const Ripemd160 Hash = "ripemd160"
const Md4 Hash = "md4"
const Md5 Hash = "md5"
const Sha1 Hash = "sha1"
const Sha256 Hash = "sha256"
const Sha384 Hash = "sha384"
const Sha3_224 Hash = "sha3-224"
const Sha3_256 Hash = "sha3-256"
const Sha3_384 Hash = "sha3-384"
const Sha3_512 Hash = "sha3-512"
const Sha512 Hash = "sha512"
const Sha512_224 Hash = "sha512-224"
const Sha512_256 Hash = "sha512-256"
const Crc32 Hash = "crc32"
const Crc64_iso Hash = "crc64-iso"
const Crc64_ecma Hash = "crc64-ecma"
const Adler32 Hash = "adler32"
const Fnv32 Hash = "fnv32"
const Fnv32a Hash = "fnv32a"
const Fnv64 Hash = "fnv64"
const Fnv64a Hash = "fnv64a"
const Fnv128 Hash = "fnv128"
const Fnv128a Hash = "fnv128a"
const Xor8 Hash = "xor8"
const Fletch16 Hash = "fletch16"
const Luhn Hash = "luhn"
const Sum16 Hash = "sum16"
const Sum32 Hash = "sum32"
const Sum64 Hash = "sum64"
const Crc8 Hash = "crc8"
const Crc16_ccitt Hash = "crc16-ccitt"

var Hashs = []Hash{
	Blake2b_256,
	Blake2b_384,
	Blake2b_512,
	Blake2s_256,
	Ripemd160,
	Md4,
	Md5,
	Sha1,
	Sha256,
	Sha384,
	Sha3_224,
	Sha3_256,
	Sha3_384,
	Sha3_512,
	Sha512,
	Sha512_224,
	Sha512_256,
	Crc32,
	Crc64_iso,
	Crc64_ecma,
	Adler32,
	Fnv32,
	Fnv32a,
	Fnv64,
	Fnv64a,
	Fnv128,
	Fnv128a,
	Xor8,
	Fletch16,
	Luhn,
	Sum16,
	Sum32,
	Sum64,
	Crc8,
	Crc16_ccitt,
}
