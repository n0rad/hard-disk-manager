package algo

type Xor8 uint8

func (d *Xor8) Reset()      { *d = 0 }
func (Xor8) Size() int      { return 1 }
func (Xor8) BlockSize() int { return 1 }

func (d *Xor8) Write(b []byte) (int, error) {
	for i := range b {
		*d ^= Xor8(b[i])
	}
	return len(b), nil
}

func (d Xor8) Sum(b []byte) []byte {
	return append(b, uint8(d))
}
