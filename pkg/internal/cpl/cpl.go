package cpl

const byteWidth = 8

type Comparable interface {
	Bytes() []byte
}

func CommonPrefixLen(x, y Comparable) int {
	bx, by := x.Bytes(), y.Bytes()

	bytes := len(bx)
	for idxOfByte := 0; idxOfByte < bytes; idxOfByte++ {
		r := bx[idxOfByte] ^ by[idxOfByte]
		if r > 0 {
			for idxOfBit := 0; idxOfBit < byteWidth; idxOfBit, r = idxOfBit+1, r<<1 {
				if r > 127 {
					return idxOfByte*byteWidth + idxOfBit
				}
			}
		}
	}

	return bytes * byteWidth
}
