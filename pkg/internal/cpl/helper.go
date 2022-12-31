package cpl

type bytesComparable []byte

func (x bytesComparable) Bytes() []byte {
	return []byte(x)
}

func FromBytes(x []byte) Comparable {
	return bytesComparable(x)
}
