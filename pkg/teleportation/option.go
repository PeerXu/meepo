package teleportation

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/transport"
)

type OFN = func(objx.Map)

type NewTeleportationSourceOption = OFN
type NewTeleportationSinkOption = OFN

func WithLogger(logger logrus.FieldLogger) OFN {
	return func(o objx.Map) {
		o["logger"] = logger
	}
}

func WithName(name string) OFN {
	return func(o objx.Map) {
		o["name"] = name
	}
}

func WithSource(addr net.Addr) OFN {
	return func(o objx.Map) {
		o["source"] = addr
	}
}

func WithSink(addr net.Addr) OFN {
	return func(o objx.Map) {
		o["sink"] = addr
	}
}

func WithTransport(t transport.Transport) OFN {
	return func(o objx.Map) {
		o["transport"] = t
	}
}

type NewDial func(network, address string) (net.Conn, error)

func WithNewDial(f NewDial) OFN {
	return func(o objx.Map) {
		o["newDial"] = f
	}
}

type DoTeleportFunc func(label string) error

func WithDoTeleportFunc(f DoTeleportFunc) OFN {
	return func(o objx.Map) {
		o["doTeleportFunc"] = f
	}
}

type OnDoTeleportFunc func() error

func WithOnDoTeleportFunc(f OnDoTeleportFunc) OFN {
	return func(o objx.Map) {
		o["onDoTeleportFunc"] = f
	}
}

type OnCloseHandler func()

func WithOnCloseHandler(h OnCloseHandler) OFN {
	return func(o objx.Map) {
		o["onCloseHandler"] = h
	}
}

type OnErrorHandler func(error)

func WithOnErrorHandler(h OnErrorHandler) OFN {
	return func(o objx.Map) {
		o["onErrorHandler"] = h
	}
}

type DialRequest struct {
	Conn net.Conn
	Quit chan struct{}
}

func NewDialRequest(conn net.Conn) *DialRequest {
	return &DialRequest{Conn: conn}
}

func NewDialRequestWithQuit(conn net.Conn, quit chan struct{}) *DialRequest {
	return &DialRequest{Conn: conn, Quit: quit}
}

func SetDialRequestChannel(c chan *DialRequest) OFN {
	return func(o objx.Map) {
		o["dialRequestChannel"] = c
	}
}
