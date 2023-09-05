package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	meepo_eventloop_core "github.com/PeerXu/meepo/pkg/meepo/eventloop/core"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
	teleportation_core "github.com/PeerXu/meepo/pkg/meepo/teleportation/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (mp *Meepo) emit(evt meepo_eventloop_interface.Event) {
	mp.eventloop.Emit(evt)
}

func (mp *Meepo) emitTransportActionNew(t transport_core.Transport) {
	dat := mp.viewToMap(ViewTransport(t))
	evt := meepo_eventloop_core.NewEvent(EVENT_TRANSPORT_ACTION_NEW, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitTransportActionClose(t transport_core.Transport) {
	dat := mp.viewToMap(ViewTransport(t))
	evt := meepo_eventloop_core.NewEvent(EVENT_TRANSPORT_ACTION_CLOSE, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitChannelActionNew(c transport_core.Channel) {
	dat := mp.viewToMap(ViewChannel(c))
	evt := meepo_eventloop_core.NewEvent(EVENT_CHANNEL_ACTION_NEW, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitChannelActionClose(c transport_core.Channel) {
	dat := mp.viewToMap(ViewChannel(c))
	evt := meepo_eventloop_core.NewEvent(EVENT_CHANNEL_ACTION_CLOSE, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitTeleportationNew(tp teleportation_core.Teleportation) {
	dat := mp.viewToMap(ViewTeleportation(tp))
	evt := meepo_eventloop_core.NewEvent(EVENT_TELEPORTATION_ACTION_NEW, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitTeleportationClose(tp teleportation_core.Teleportation) {
	dat := mp.viewToMap(ViewTeleportation(tp))
	evt := meepo_eventloop_core.NewEvent(EVENT_TELEPORTATION_ACTION_CLOSE, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitTransportStateChange(t transport_core.Transport) {
	dat := mp.viewToMap(ViewTransport(t))
	name := ConvertTransportStateToEventName(t.State())
	evt := meepo_eventloop_core.NewEvent(name, dat)
	mp.emit(evt)
}

func (mp *Meepo) emitChannelStateChange(x addr.Addr, c transport_core.Channel) {
	dat := mp.viewToMap(ViewChannelWithAddr(c, x))
	name := ConvertChannelStateToEventName(c.State())
	evt := meepo_eventloop_core.NewEvent(name, dat)
	mp.emit(evt)
}
