package meepo_core

import (
	"context"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func wrapHandleFunc[T any](name string, h rpc_interface.Handler, fn func(context.Context, any) (any, error)) {
	h.Handle(name, rpc_core.WrapHandleFuncGenerics[T](fn))
}

func (mp *Meepo) AsTrackerdHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")

	wrapHandleFunc[crypto_core.Packet](tracker_core.METHOD_NEW_TRANSPORT, h, mp.hdrOnNewTransport)
	wrapHandleFunc[GetCandidatesRequest](tracker_core.METHOD_GET_CANDIDATES, h, mp.hdrOnGetCandidates)
	wrapHandleFunc[crypto_core.Packet](tracker_core.METHOD_ADD_PEER_CONNECTION, h, mp.hdrOnAddPeerConnection)

	return h
}

func (mp *Meepo) AsAPIHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")

	// system
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_GET_VERSION, h, mp.hdrAPIGetVersion)
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_WHOAMI, h, mp.hdrAPIWhoami)
	wrapHandleFunc[sdk_interface.PingRequest](sdk_core.METHOD_PING, h, mp.hdrAPIPing)
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_DIAGNOSTIC, h, mp.hdrAPIDiagnostic)

	h.HandleStream("watchEvents", mp.hdrStreamAPIWatchEvents)

	// teleportation
	wrapHandleFunc[sdk_interface.NewTeleportationRequest](sdk_core.METHOD_NEW_TELEPORTATION, h, mp.hdrAPINewTeleportation)
	wrapHandleFunc[sdk_interface.CloseTeleportationRequest](sdk_core.METHOD_CLOSE_TELEPORTATION, h, mp.hdrAPICloseTeleportation)
	wrapHandleFunc[sdk_interface.GetTeleportationRequest](sdk_core.METHOD_GET_TELEPORTATION, h, mp.hdrAPIGetTeleportation)
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_LIST_TELEPORTATIONS, h, mp.hdrAPIListTeleportations)
	wrapHandleFunc[sdk_interface.TeleportRequest](sdk_core.METHOD_TELEPORT, h, mp.hdrAPITeleport)

	// transport
	wrapHandleFunc[sdk_interface.NewTransportRequest](sdk_core.METHOD_NEW_TRANSPORT, h, mp.hdrAPINewTransport)
	wrapHandleFunc[sdk_interface.CloseTransportRequest](sdk_core.METHOD_CLOSE_TRANSPORT, h, mp.hdrAPICloseTransport)
	wrapHandleFunc[sdk_interface.GetTransportRequest](sdk_core.METHOD_GET_TRANSPORT, h, mp.hdrAPIGetTransport)
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_LIST_TRANSPORTS, h, mp.hdrAPIListTransports)

	// channel
	wrapHandleFunc[sdk_interface.CloseChannelRequest](sdk_core.METHOD_CLOSE_CHANNEL, h, mp.hdrAPICloseChannel)
	wrapHandleFunc[sdk_interface.GetChannelRequest](sdk_core.METHOD_GET_CHANNEL, h, mp.hdrAPIGetChannel)
	wrapHandleFunc[rpc_core.EMPTY](sdk_core.METHOD_LIST_CHANNELS, h, mp.hdrAPIListChannels)
	wrapHandleFunc[sdk_interface.ListChannelsByTarget](sdk_core.METHOD_LIST_CHANNELS_BY_TARGET, h, mp.hdrAPIListChannelsByTarget)

	return h
}
