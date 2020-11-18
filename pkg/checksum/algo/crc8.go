package algo

type Crc8 uint8

func (c *Crc8) Reset()      { *c = 0 }
func (Crc8) Size() int      { return 1 }
func (Crc8) BlockSize() int { return 1 }

func (c Crc8) calc(b byte) Crc8 {
	d := c ^ Crc8(b)
	for i := 0; i < 8; i++ {
		if d&0x80 != 0 {
			d <<= 1
			d ^= 0x07
		} else {
			d <<= 1
		}
	}
	return d
}

func (c *Crc8) Write(b []byte) (int, error) {
	for i := range b {
		*c = c.calc(b[i])
	}
	return len(b), nil
}

func (c Crc8) Sum(b []byte) []byte {
	return append(b, byte(c))
}
