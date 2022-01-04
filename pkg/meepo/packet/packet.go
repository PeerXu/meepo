package packet

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/util/msgpack"
)

type Packet interface {
	Header() Header
	Err() error
	Data(out interface{}) error
	SetSignature([]byte) Packet
	UnsetSignature() Packet

	msgpack.Marshaler
	msgpack.Unmarshaler
}

type NewPacketOption = ofn.OFN

func WithError(err error) ofn.OFN {
	return func(o ofn.Option) {
		o["error"] = err
	}
}

func WithData(v interface{}) ofn.OFN {
	return func(o ofn.Option) {
		o["data"] = v
	}
}

func NewPacket(h Header, opts ...NewPacketOption) (Packet, error) {
	p := &packet{
		header: h,
	}

	o := ofn.NewOption(map[string]interface{}{})

	for _, opt := range opts {
		opt(o)
	}

	if v := o.Get("error").Inter(); v != nil {
		p.err = v.(error)
		return p, nil
	}

	if v := o.Get("data").Inter(); v != nil {
		raw, err := msgpack.Marshal(v)
		if err != nil {
			return nil, err
		}

		p.rawData = raw
		return p, nil
	}

	return p, nil
}

type packet struct {
	header  Header
	err     error
	rawData []byte
}

func (p *packet) Header() Header {
	return p.header
}

func (p *packet) Err() error {
	return p.err
}

func (p *packet) Data(out interface{}) error {
	return msgpack.Unmarshal(p.rawData, &out)
}

func (p *packet) SetSignature(b []byte) Packet {
	pp := *p
	pp.header = p.Header().SetSignature(b)
	return &pp
}

func (p *packet) UnsetSignature() Packet {
	pp := *p
	pp.header = p.Header().UnsetSignature()
	return &pp
}

var _ msgpack.Marshaler = (*packet)(nil)

func (p *packet) MarshalMsgpack() ([]byte, error) {
	bHdr, err := msgpack.Marshal(p.header)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		"header": bHdr,
	}

	if p.err != nil {
		m["error"] = p.err.Error()
	} else {
		m["data"] = base64.StdEncoding.EncodeToString(p.rawData)
	}

	return msgpack.Marshal(m)
}

var _ msgpack.Unmarshaler = (*packet)(nil)

func (p *packet) UnmarshalMsgpack(b []byte) error {
	var m map[string]interface{}

	err := msgpack.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	iHdr, ok := m["header"]
	if !ok {
		return ErrInvalidPacket
	}

	if p.header, err = UnmarshalHeader(iHdr.([]byte)); err != nil {
		return err
	}

	iErr, ok := m["error"]
	if ok {
		p.err = fmt.Errorf(iErr.(string))
		return nil
	}

	iData, ok := m["data"]
	if !ok {
		return ErrInvalidPacket
	}

	p.rawData, _ = base64.StdEncoding.DecodeString(cast.ToString(iData))

	return nil
}

func MarshalPacket(p Packet) ([]byte, error) {
	return msgpack.Marshal(p)
}

func UnmarshalPacket(b []byte) (Packet, error) {
	var p packet
	if err := msgpack.Unmarshal(b, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
