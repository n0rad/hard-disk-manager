package algo

type Sum16 uint16

func (d *Sum16) Reset()      { *d = 0 }
func (Sum16) Size() int      { return 2 }
func (Sum16) BlockSize() int { return 1 }

func (d *Sum16) Write(b []byte) (int, error) {
	for i := range b {
		*d += Sum16(b[i])
	}
	return len(b), nil
}

func (d Sum16) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8))
}
