package routing_table

import "bytes"

const (
	BYTE_WIDTH = 8
)

var (
	_0XFF = uint8(0xff)
)

type id []byte

func FromBytes(x []byte) ID {
	return id(x)
}

func (x id) Bytes() []byte {
	return x
}

func (x id) Equal(y ID) bool {
	return bytes.Equal(x.Bytes(), y.Bytes())
}
