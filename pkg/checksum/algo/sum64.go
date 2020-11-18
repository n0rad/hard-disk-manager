package algo

type Sum64 uint64

func (d *Sum64) Reset()      { *d = 0 }
func (Sum64) Size() int      { return 8 }
func (Sum64) BlockSize() int { return 1 }

func (d *Sum64) Write(b []byte) (int, error) {
	for i := range b {
		*d += Sum64(b[i])
	}
	return len(b), nil
}

func (d Sum64) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8), byte(d>>16), byte(d>>24), byte(d>>32), byte(d>>40), byte(d>>48), byte(d>>56))
}
