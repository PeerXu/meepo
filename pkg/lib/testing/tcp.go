package testing

import (
	"context"
	"io"
	"net"
)

type EchoServer struct {
	Listener net.Listener
}

func (es *EchoServer) Serve(ctx context.Context) (err error) {
	defer es.Listener.Close()
	for {
		conn, err := es.Listener.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			io.Copy(conn, conn) // nolint
		}(conn)
	}
}

func (es *EchoServer) Terminate(ctx context.Context) error {
	if es.Listener != nil {
		if err := es.Listener.Close(); err != nil {
			return err
		}
	}

	es.Listener = nil

	return nil
}

func NewEchoServer() (*EchoServer, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	return &EchoServer{Listener: lis}, nil
}
