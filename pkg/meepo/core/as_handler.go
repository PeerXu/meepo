package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) AsTrackerdHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")
	h.Handle("newTransport", rpc_core.WrapHandleFunc(mp.newOnNewTransportRequest, mp.hdrOnNewTransport))
	h.Handle("getCandidates", rpc_core.WrapHandleFunc(mp.newOnGetCandidatesRequest, mp.hdrOnGetCandidates))
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
		{"getVersion", rpc_core.NO_REQUEST, mp.hdrAPIGetVersion},
		{"whoami", rpc_core.NO_REQUEST, mp.hdrAPIWhoami},
		{"ping", func() any { return &sdk_interface.PingRequest{} }, mp.hdrAPIPing},
		{"diagnostic", rpc_core.NO_REQUEST, mp.hdrAPIDiagnostic},

		// teleportation
		{"newTeleportation", func() any { return &sdk_interface.NewTeleportationRequest{} }, mp.hdrAPINewTeleportation},
		{"closeTeleportation", func() any { return &sdk_interface.CloseTeleportationRequest{} }, mp.hdrAPICloseTeleportation},
		{"getTeleportation", func() any { return &sdk_interface.GetTeleportationRequest{} }, mp.hdrAPIGetTeleportation},
		{"listTeleportations", rpc_core.NO_REQUEST, mp.hdrAPIListTeleportations},
		{"teleport", func() any { return &sdk_interface.TeleportRequest{} }, mp.hdrAPITeleport},

		// transport
		{"newTransport", func() any { return &sdk_interface.NewTransportRequest{} }, mp.hdrAPINewTransport},
		{"closeTransport", func() any { return &sdk_interface.CloseTransportRequest{} }, mp.hdrAPICloseTransport},
		{"getTransport", func() any { return &sdk_interface.GetTransportRequest{} }, mp.hdrAPIGetTransport},
		{"listTransports", rpc_core.NO_REQUEST, mp.hdrAPIListTransports},

		// channel
		{"closeChannel", func() any { return &sdk_interface.CloseChannelRequest{} }, mp.hdrAPICloseChannel},
		{"getChannel", func() any { return &sdk_interface.GetChannelRequest{} }, mp.hdrAPIGetChannel},
		{"listChannels", rpc_core.NO_REQUEST, mp.hdrAPIListChannels},
		{"listChannelsByTarget", func() any { return &sdk_interface.ListChannelsByTarget{} }, mp.hdrAPIListChannelsByTarget},
	} {
		h.Handle(s.name, rpc_core.WrapHandleFunc(s.newRequest, s.fn))
	}

	return h
}
