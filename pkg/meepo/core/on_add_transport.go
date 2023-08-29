package meepo_core

import (
	"context"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func wrapTransportHandleFunc[T any](name string, t Transport, fn func(context.Context, any) (any, error)) {
	t.Handle(name, transport_core.WrapHandleFuncGenerics[T](fn))
}

func (mp *Meepo) onAddWebrtcTransportNL(t Transport) {
	wrapTransportHandleFunc[crypto_core.Packet](tracker_core.METHOD_NEW_TRANSPORT, t, mp.hdrOnNewTransport)
	wrapTransportHandleFunc[GetCandidatesRequest](tracker_core.METHOD_GET_CANDIDATES, t, mp.hdrOnGetCandidates)
	wrapTransportHandleFunc[crypto_core.Packet](tracker_core.METHOD_ADD_PEER_CONNECTION, t, mp.hdrOnAddPeerConnection)

	wrapTransportHandleFunc[PingRequest](METHOD_PING, t, mp.hdrOnPing)
	wrapTransportHandleFunc[PermitRequest](METHOD_PERMIT, t, mp.hdrOnPermit)

	mp.transports[t.Addr()] = t
}

func (mp *Meepo) onAddPipeTransportNL(t Transport) {
	wrapTransportHandleFunc[PingRequest](METHOD_PING, t, mp.hdrOnPing)
	wrapTransportHandleFunc[PermitRequest](METHOD_PERMIT, t, mp.hdrOnPermit)

	mp.transports[t.Addr()] = t
}
