package listenerer_http

import (
	"io"
	"net"
	"net/http"
	"sync"
)

type HttpConnectConn struct {
	reader     io.Reader
	writer     io.Writer
	request    *http.Request
	close      chan struct{}
	closeOnce  sync.Once
	enable     chan struct{}
	enableOnce sync.Once
}

type HttpGetConn struct {
	rd1, rd2   *io.PipeReader
	wr1, wr2   *io.PipeWriter
	remoteAddr net.Addr
	close      chan struct{}
	closeOnce  sync.Once
	enable     chan struct{}
	enableOnce sync.Once
}

func NewHttpGetConn(network, address string) *HttpGetConn {
	rd1, wr1 := io.Pipe()
	rd2, wr2 := io.Pipe()
	ra := &net.UnixAddr{Net: network, Name: address}
	return &HttpGetConn{
		rd1:        rd1,
		rd2:        rd2,
		wr1:        wr1,
		wr2:        wr2,
		remoteAddr: ra,
		close:      make(chan struct{}),
		enable:     make(chan struct{}),
	}
}

type HttpGetPipeConn struct {
	*HttpGetConn
	*io.PipeReader
	*io.PipeWriter
}
