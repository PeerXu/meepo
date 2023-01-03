package dialer_interface

import "io"

type Conn interface {
	io.ReadWriteCloser
}
