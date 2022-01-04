package packet

import (
	"encoding/base64"

	"github.com/PeerXu/meepo/pkg/ofn"
)

type BroadcastPacket interface {
	Header() BroadcastHeader
	Err() error
	Packet() Packet
	SetPacket(Packet) BroadcastPacket
}

type broadcastPacket struct {
	header BroadcastHeader
	packet Packet
}

func (p *broadcastPacket) Header() BroadcastHeader {
	return p.header
}

func (p *broadcastPacket) Err() error {
	if p.packet == nil {
		return ErrPacketIsNil
	}

	return p.packet.Err()
}

func (p *broadcastPacket) Packet() Packet {
	return p.packet
}

func (p *broadcastPacket) SetPacket(t Packet) BroadcastPacket {
	pp := *p
	pp.packet = t
	return &pp
}

type NewBroadcastPacketOption = ofn.OFN

func WithPacket(p Packet) ofn.OFN {
	return func(o ofn.Option) {
		o["packet"] = p
	}
}

func NewBroadcastPacket(h BroadcastHeader, opts ...NewBroadcastPacketOption) (BroadcastPacket, error) {
	o := ofn.NewOption(map[string]interface{}{})
	for _, opt := range opts {
		opt(o)
	}

	p, ok := o.Get("packet").Inter().(Packet)
	if !ok {
		return nil, ErrPacketIsNil
	}

	return &broadcastPacket{
		header: h,
		packet: p,
	}, nil
}

type _broadcastPacket struct {
	Hop    int32
	Packet string
}

func PacketToBroadcastPacketE(p Packet) (bp BroadcastPacket, err error) {
	var t _broadcastPacket
	if err = p.Data(&t); err != nil {
		return
	}

	pbuf, err := base64.StdEncoding.DecodeString(t.Packet)
	if err != nil {
		return
	}

	np, err := UnmarshalPacket(pbuf)
	if err != nil {
		return
	}

	bhdr := NewBroadcastHeader(
		p.Header().Session(),
		p.Header().Source(),
		p.Header().Destination(),
		p.Header().Type(),
		p.Header().Method(),
		t.Hop,
	)

	return NewBroadcastPacket(bhdr, WithPacket(np))
}

func BroadcastPacketToPacketE(bp BroadcastPacket) (p Packet, err error) {
	var pbuf []byte

	bhdr := bp.Header()

	pbuf, err = MarshalPacket(bp.Packet())
	if err != nil {
		return
	}

	hdr := NewHeader(bhdr.Session(), bhdr.Source(), bhdr.Destination(), bhdr.Type(), bhdr.Method())
	p, err = NewPacket(hdr, WithData(&_broadcastPacket{
		Hop:    bhdr.Hop(),
		Packet: base64.StdEncoding.EncodeToString(pbuf),
	}))

	return
}

func BroadcastPacketToPacket(bp BroadcastPacket) (p Packet) {
	p, _ = BroadcastPacketToPacketE(bp)
	return
}
