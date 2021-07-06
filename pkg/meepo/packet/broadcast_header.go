package packet

import "fmt"

type BroadcastHeader interface {
	Session() int32
	Source() string
	Destination() string
	Type() Type
	Method() string
	Hop() int32
}

func NewBroadcastHeader(sess int32, src, dst string, typ Type, meth string, hop int32) BroadcastHeader {
	return &broadcastHeader{
		session:     sess,
		source:      src,
		destination: dst,
		typ:         typ,
		method:      meth,
		hop:         hop,
	}
}

func InvertBroadcastHeader(in BroadcastHeader) (out BroadcastHeader) {
	var typ Type
	switch in.Type() {
	case BroadcastRequest:
		typ = BroadcastResponse
	case BroadcastResponse:
		typ = BroadcastRequest
	default:
		panic(fmt.Errorf("unexpected broadcast type: %v", in.Type()))
	}

	return &broadcastHeader{
		session:     in.Session(),
		source:      in.Destination(),
		destination: in.Source(),
		typ:         typ,
		method:      in.Method(),
		hop:         in.Hop(),
	}
}

type broadcastHeader struct {
	session     int32
	source      string
	destination string
	typ         Type
	method      string
	hop         int32
}

func (h *broadcastHeader) Session() int32 {
	return h.session
}

func (h *broadcastHeader) Source() string {
	return h.source
}

func (h *broadcastHeader) Destination() string {
	return h.destination
}

func (h *broadcastHeader) Type() Type {
	return h.typ
}

func (h *broadcastHeader) Method() string {
	return h.method
}

func (h *broadcastHeader) Hop() int32 {
	return h.hop
}
