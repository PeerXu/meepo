package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func (mp *Meepo) AsTrackerdHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")
	h.Handle(tracker_core.METHOD_NEW_TRANSPORT, rpc_core.WrapHandleFunc(mp.newOnNewTransportRequest, mp.hdrOnNewTransport))
	h.Handle(tracker_core.METHOD_GET_CANDIDATES, rpc_core.WrapHandleFunc(mp.newOnGetCandidatesRequest, mp.hdrOnGetCandidates))
	h.Handle(tracker_core.METHOD_ADD_PEER_CONNECTION, rpc_core.WrapHandleFunc(mp.newOnAddPeerConnectionRequest, mp.hdrOnAddPeerConnection))
	return h
}

func (mp *Meepo) AsAPIHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")

	for _, s := range []struct {
		name       string
		newRequest func() any
		fn         func(context.Context, any) (any, error)
	}{
		// system
		{sdk_core.METHOD_GET_VERSION, rpc_core.NO_REQUEST, mp.hdrAPIGetVersion},
		{sdk_core.METHOD_WHOAMI, rpc_core.NO_REQUEST, mp.hdrAPIWhoami},
		{sdk_core.METHOD_PING, func() any { return &sdk_interface.PingRequest{} }, mp.hdrAPIPing},
		{sdk_core.METHOD_DIAGNOSTIC, rpc_core.NO_REQUEST, mp.hdrAPIDiagnostic},

		// teleportation
		{sdk_core.METHOD_NEW_TELEPORTATION, func() any { return &sdk_interface.NewTeleportationRequest{} }, mp.hdrAPINewTeleportation},
		{sdk_core.METHOD_CLOSE_TELEPORTATION, func() any { return &sdk_interface.CloseTeleportationRequest{} }, mp.hdrAPICloseTeleportation},
		{sdk_core.METHOD_GET_TELEPORTATION, func() any { return &sdk_interface.GetTeleportationRequest{} }, mp.hdrAPIGetTeleportation},
		{sdk_core.METHOD_LIST_TELEPORTATIONS, rpc_core.NO_REQUEST, mp.hdrAPIListTeleportations},
		{sdk_core.METHOD_TELEPORT, func() any { return &sdk_interface.TeleportRequest{} }, mp.hdrAPITeleport},

		// transport
		{sdk_core.METHOD_NEW_TRANSPORT, func() any { return &sdk_interface.NewTransportRequest{} }, mp.hdrAPINewTransport},
		{sdk_core.METHOD_CLOSE_TRANSPORT, func() any { return &sdk_interface.CloseTransportRequest{} }, mp.hdrAPICloseTransport},
		{sdk_core.METHOD_GET_TRANSPORT, func() any { return &sdk_interface.GetTransportRequest{} }, mp.hdrAPIGetTransport},
		{sdk_core.METHOD_LIST_TRANSPORTS, rpc_core.NO_REQUEST, mp.hdrAPIListTransports},

		// channel
		{sdk_core.METHOD_CLOSE_CHANNEL, func() any { return &sdk_interface.CloseChannelRequest{} }, mp.hdrAPICloseChannel},
		{sdk_core.METHOD_GET_CHANNEL, func() any { return &sdk_interface.GetChannelRequest{} }, mp.hdrAPIGetChannel},
		{sdk_core.METHOD_LIST_CHANNELS, rpc_core.NO_REQUEST, mp.hdrAPIListChannels},
		{sdk_core.METHOD_LIST_CHANNELS_BY_TARGET, func() any { return &sdk_interface.ListChannelsByTarget{} }, mp.hdrAPIListChannelsByTarget},
	} {
		h.Handle(s.name, rpc_core.WrapHandleFunc(s.newRequest, s.fn))
	}

	return h
}
