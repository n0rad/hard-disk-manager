package algo

type Luhn uint64

func (d *Luhn) Reset()      { *d = 0 }
func (Luhn) Size() int      { return 8 }
func (Luhn) BlockSize() int { return 1 }

func (d *Luhn) Write(b []byte) (int, error) {
	tab := [...]uint64{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
	odd := len(b) & 1

	var sum uint64
	for i, c := range b {
		if c < '0' || c > '9' {
			c %= 10
		} else {
			c -= '0'
		}
		if i&1 == odd {
			sum += tab[c]
		} else {
			sum += uint64(c)
		}
	}
	*d = Luhn(sum)

	return len(b), nil
}

func (d Luhn) Sum(b []byte) []byte {
	return append(b, byte(d), byte(d>>8), byte(d>>16), byte(d>>24),
		byte(d>>32), byte(d>>40), byte(d>>48), byte(d>>56))
}
