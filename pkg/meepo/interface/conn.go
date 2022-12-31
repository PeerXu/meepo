package meepo_interface

import (
	"io"
)

type Conn interface {
	io.ReadWriteCloser
}
