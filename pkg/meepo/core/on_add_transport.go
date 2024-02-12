package meepo_core

import (
	"context"

	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func wrapTransportHandleFunc[IT, OT any](name string, t Transport, fn func(context.Context, IT) (OT, error)) {
	t.Handle(name, transport_core.WrapHandleFuncGenerics(fn))
}

func (mp *Meepo) onAddWebrtcTransportNL(t Transport) {
	wrapTransportHandleFunc(tracker_core.METHOD_NEW_TRANSPORT, t, mp.hdrOnNewTransport)
	wrapTransportHandleFunc(tracker_core.METHOD_GET_CANDIDATES, t,
		RequireProtocolVersion[GetCandidatesRequest, GetCandidatesResponse]("v0.1.0", "")(mp.onGetCandidates),
	)
	wrapTransportHandleFunc(tracker_core.METHOD_ADD_PEER_CONNECTION, t, mp.hdrOnAddPeerConnection)

	wrapTransportHandleFunc(METHOD_PING, t, mp.onPing)
	wrapTransportHandleFunc(METHOD_PERMIT, t, mp.onPermit)

	mp.transports[t.Addr()] = t
}

func (mp *Meepo) onAddPipeTransportNL(t Transport) {
	wrapTransportHandleFunc(METHOD_PING, t, mp.onPing)
	wrapTransportHandleFunc(METHOD_PERMIT, t, mp.onPermit)

	mp.transports[t.Addr()] = t
}
