package packet

import (
	"encoding/base64"

	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/util/msgpack"
)

type Header interface {
	Session() int32
	Source() string
	Destination() string
	Type() Type
	Method() string
	Signature() []byte

	SetSignature([]byte) Header
	UnsetSignature() Header

	msgpack.Marshaler
	msgpack.Unmarshaler
}

func NewHeader(sess int32, src, dst string, typ Type, meth string) Header {
	return &header{
		session:     sess,
		source:      src,
		destination: dst,
		typ:         typ,
		method:      meth,
	}
}

type header struct {
	session     int32
	source      string
	destination string
	typ         Type
	method      string
	signature   []byte
}

func (h *header) Session() int32 {
	return h.session
}

func (h *header) Source() string {
	return h.source
}

func (h *header) Destination() string {
	return h.destination
}

func (h *header) Type() Type {
	return h.typ
}

func (h *header) Method() string {
	return h.method
}

func (h *header) Signature() []byte {
	return h.signature
}

var _ msgpack.Marshaler = (*header)(nil)

func (h *header) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(map[string]interface{}{
		"session":     h.session,
		"source":      h.source,
		"destination": h.destination,
		"type":        h.typ,
		"method":      h.method,
		"signature":   base64.StdEncoding.EncodeToString(h.signature),
	})
}

var _ msgpack.Unmarshaler = (*header)(nil)

func (h *header) UnmarshalMsgpack(b []byte) error {
	var m map[string]interface{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return err
	}

	h.session = cast.ToInt32(m["session"])
	h.source = cast.ToString(m["source"])
	h.destination = cast.ToString(m["destination"])
	h.typ = Type(cast.ToString(m["type"]))
	h.method = cast.ToString(m["method"])
	h.signature, _ = base64.StdEncoding.DecodeString(cast.ToString(m["signature"]))

	return nil
}

func (h *header) SetSignature(b []byte) Header {
	hh := *h
	hh.signature = b
	return &hh
}

func (h *header) UnsetSignature() Header {
	hh := *h
	hh.signature = nil
	return &hh
}

func MarshalHeader(h Header) ([]byte, error) {
	return msgpack.Marshal(h)
}

func UnmarshalHeader(b []byte) (Header, error) {
	var h header
	if err := msgpack.Unmarshal(b, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

func InvertHeader(in Header) Header {
	var typ Type
	switch in.Type() {
	case Request:
		typ = Response
	case Response:
		typ = Request
	case BroadcastRequest:
		typ = BroadcastResponse
	case BroadcastResponse:
		typ = BroadcastRequest
	}

	return &header{
		session:     in.Session(),
		source:      in.Destination(),
		destination: in.Source(),
		typ:         typ,
		method:      in.Method(),
	}
}
