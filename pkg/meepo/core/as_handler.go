package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func wrapHandleFunc[IT any, OT any](name string, h rpc_interface.Handler, fn func(context.Context, IT) (OT, error)) {
	h.Handle(name, rpc_core.WrapHandleFuncGenerics(fn))
}

func (mp *Meepo) AsTrackerdHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")

	wrapHandleFunc(tracker_core.METHOD_NEW_TRANSPORT, h, mp.hdrOnNewTransport)
	wrapHandleFunc(tracker_core.METHOD_GET_CANDIDATES, h, mp.onGetCandidates)
	wrapHandleFunc(tracker_core.METHOD_ADD_PEER_CONNECTION, h, mp.hdrOnAddPeerConnection)

	return h
}

func (mp *Meepo) AsAPIHandler() rpc_interface.Handler {
	h, _ := rpc.NewHandler("default")

	// system
	wrapHandleFunc(sdk_core.METHOD_GET_VERSION, h, mp.apiGetVersion)
	wrapHandleFunc(sdk_core.METHOD_WHOAMI, h, mp.apiWhoami)
	wrapHandleFunc(sdk_core.METHOD_PING, h, mp.apiPing)
	wrapHandleFunc(sdk_core.METHOD_DIAGNOSTIC, h, mp.apiDiagnostic)

	h.HandleStream(sdk_core.STREAM_METHOD_WATCH_EVENTS, mp.hdrStreamAPIWatchEvents)

	// teleportation
	wrapHandleFunc(sdk_core.METHOD_NEW_TELEPORTATION, h, mp.apiNewTeleportation)
	wrapHandleFunc(sdk_core.METHOD_CLOSE_TELEPORTATION, h, mp.apiCloseTeleportation)
	wrapHandleFunc(sdk_core.METHOD_GET_TELEPORTATION, h, mp.apiGetTeleportation)
	wrapHandleFunc(sdk_core.METHOD_LIST_TELEPORTATIONS, h, mp.apiListTeleportations)
	wrapHandleFunc(sdk_core.METHOD_TELEPORT, h, mp.apiTeleport)

	// transport
	wrapHandleFunc(sdk_core.METHOD_NEW_TRANSPORT, h, mp.apiNewTransport)
	wrapHandleFunc(sdk_core.METHOD_CLOSE_TRANSPORT, h, mp.apiCloseTransport)
	wrapHandleFunc(sdk_core.METHOD_GET_TRANSPORT, h, mp.apiGetTransport)
	wrapHandleFunc(sdk_core.METHOD_LIST_TRANSPORTS, h, mp.apiListTransports)

	// channel
	wrapHandleFunc(sdk_core.METHOD_CLOSE_CHANNEL, h, mp.apiCloseChannel)
	wrapHandleFunc(sdk_core.METHOD_GET_CHANNEL, h, mp.apiGetChannel)
	wrapHandleFunc(sdk_core.METHOD_LIST_CHANNELS, h, mp.apiListChannels)
	wrapHandleFunc(sdk_core.METHOD_LIST_CHANNELS_BY_TARGET, h, mp.apiListChannelsByTarget)

	return h
}
