package algo

type Fletch16 struct {
	sum1, sum2 uint16
}

func (d *Fletch16) Reset()      { d.sum1, d.sum2 = 0, 0 }
func (Fletch16) Size() int      { return 2 }
func (Fletch16) BlockSize() int { return 1 }

func (d *Fletch16) Write(b []byte) (int, error) {
	for i := range b {
		d.sum1 = (d.sum1 + uint16(b[i])) % 255
		d.sum2 = (d.sum2 + d.sum1) % 255
	}
	return len(b), nil
}

func (d *Fletch16) Sum(b []byte) []byte {
	return append(b, byte(d.sum2), byte(d.sum1))
}
