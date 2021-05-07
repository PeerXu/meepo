package conn

import (
	"io"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
)

type rwConn struct {
	reader io.Reader
	writer io.Writer
	local  net.Addr
	remote net.Addr
}

func (c *rwConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *rwConn) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

func (c *rwConn) Close() error {
	var eg errgroup.Group

	rdCloser, ok := c.reader.(io.Closer)
	if ok {
		eg.Go(rdCloser.Close)
	}

	wrCloser, ok := c.writer.(io.Closer)
	if ok {
		eg.Go(wrCloser.Close)
	}

	return eg.Wait()
}

func (c *rwConn) LocalAddr() net.Addr {
	return c.local
}

func (c *rwConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *rwConn) SetDeadline(t time.Time) error {
	panic("unimplemented")
}

func (c *rwConn) SetReadDeadline(t time.Time) error {
	panic("unimplemented")
}

func (c *rwConn) SetWriteDeadline(t time.Time) error {
	panic("unimplemented")
}

func NewRWConn(reader io.Reader, writer io.Writer, local net.Addr, remote net.Addr) *rwConn {
	return &rwConn{
		reader: reader,
		writer: writer,
		local:  local,
		remote: remote,
	}
}
