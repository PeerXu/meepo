package meepo_core

import (
	"context"
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/routing_table"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func Addr2ID(x addr.Addr) routing_table.ID {
	return routing_table.FromBytes(x.Bytes())
}

func ID2Addr(x routing_table.ID) addr.Addr {
	return addr.Must(addr.FromBytesWithoutMagicCode(x.Bytes()))
}

func IDs2Addrs(xs []routing_table.ID) (ys []addr.Addr) {
	for _, x := range xs {
		ys = append(ys, ID2Addr(x))
	}
	return
}

func Addrs2IDs(xs []Addr) (ys []routing_table.ID) {
	for _, x := range xs {
		ys = append(ys, Addr2ID(x))
	}
	return
}

func ContainAddr(xs []Addr, t Addr) bool {
	for _, x := range xs {
		if x.Equal(t) {
			return true
		}
	}
	return false
}

func ViewTransport(t Transport) sdk_interface.TransportView {
	return sdk_interface.TransportView{
		Addr:  t.Addr().String(),
		State: t.State().String(),
	}
}

func ViewTransports(ts []Transport) (tvs []sdk_interface.TransportView) {
	for _, t := range ts {
		tvs = append(tvs, ViewTransport(t))
	}
	return
}

func ViewChannel(c Channel) sdk_interface.ChannelView {
	return sdk_interface.ChannelView{
		ID:          c.ID(),
		Mode:        c.Mode(),
		State:       c.State().String(),
		IsSource:    c.IsSource(),
		IsSink:      c.IsSink(),
		SinkNetwork: c.SinkAddr().Network(),
		SinkAddress: c.SinkAddr().String(),
	}
}

func ViewChannelWithAddr(c Channel, target addr.Addr) sdk_interface.ChannelView {
	cv := ViewChannel(c)
	cv.Addr = target.String()
	return cv
}

func ViewChannelsWithAddr(cs []Channel, target addr.Addr) (cvs []sdk_interface.ChannelView) {
	for _, c := range cs {
		cvs = append(cvs, ViewChannelWithAddr(c, target))
	}
	return
}

func ViewTeleportation(tp Teleportation) sdk_interface.TeleportationView {
	return sdk_interface.TeleportationView{
		ID:            tp.ID(),
		Mode:          tp.Mode(),
		Addr:          tp.Addr().String(),
		SourceNetwork: tp.SourceAddr().Network(),
		SourceAddress: tp.SourceAddr().String(),
		SinkNetwork:   tp.SinkAddr().Network(),
		SinkAddress:   tp.SinkAddr().String(),
	}
}

func ViewTeleportations(tps []Teleportation) (tpvs []sdk_interface.TeleportationView) {
	for _, tp := range tps {
		tpvs = append(tpvs, ViewTeleportation(tp))
	}
	return
}

func (mp *Meepo) newTeleportationID() string {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	for {
		id := fmt.Sprintf("%08x", mp.randSrc.Int63()&0xffffffff)
		_, found := mp.teleportations[id]
		if !found {
			return id
		}
	}
}

func (mp *Meepo) newLabel(ns string) string {
	return fmt.Sprintf("%s#%016x", ns, mp.randSrc.Int63())
}

// nolint:unused
func (mp *Meepo) hdrAPIUnimplemented(context.Context, any) (any, error) { panic("unimplemented") }
