package algo

type Crc16ccitt uint16

func (c *Crc16ccitt) Reset()      { *c = 0xffff }
func (Crc16ccitt) Size() int      { return 2 }
func (Crc16ccitt) BlockSize() int { return 1 }

func (c *Crc16ccitt) Write(b []byte) (int, error) {
	for i := range b {
		x := *c>>8 ^ Crc16ccitt(b[i])
		x ^= x >> 4
		*c = (*c << 8) ^ Crc16ccitt((x<<12)^(x<<5)^x)
	}
	return len(b), nil
}

func (c Crc16ccitt) Sum(b []byte) []byte {
	return append(b, uint8(c>>8), uint8(c))
}
