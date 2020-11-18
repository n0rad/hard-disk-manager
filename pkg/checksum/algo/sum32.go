package algo

type Sum32 uint32

func (d *Sum32) Reset()      { *d = 0 }
func (Sum32) Size() int      { return 4 }
func (Sum32) BlockSize() int { return 1 }

func (d *Sum32) Write(b []byte) (int, error) {
	for i := range b {
		*d += Sum32(b[i])
	}
	return len(b), nil
}

func (d Sum32) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8), byte(d>>16), byte(d>>24))
}
