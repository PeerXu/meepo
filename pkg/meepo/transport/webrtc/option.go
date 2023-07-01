package transport_webrtc

import (
	"time"

	"github.com/pion/webrtc/v3"

	C "github.com/PeerXu/meepo/pkg/lib/constant"
	"github.com/PeerXu/meepo/pkg/lib/option"
	mrand "github.com/PeerXu/meepo/pkg/lib/rand"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type GatherFunc func(Session, webrtc.SessionDescription) (webrtc.SessionDescription, error)
type GatherDoneFunc func(Session, webrtc.SessionDescription, error)
type NewPeerConnectionFunc func() (*webrtc.PeerConnection, error)

const (
	OPTION_ROLE                      = "role"
	OPTION_TEMP_DATA_CHANNEL_TIMEOUT = "tempDataChannelTimeout"
	OPTION_OFFER                     = "offer"
	OPTION_ANSWER                    = "answer"
	OPTION_GATHER_TIMEOUT            = "gatherTimeout"
	OPTION_GATHER_ON_NEW_FUNC        = "gatherOnNewFunc"
	OPTION_GATHER_DONE_ON_NEW_FUNC   = "gatherDoneOnNewFunc"
	OPTION_GATHER_FUNC               = "gatherFunc"

	OPTION_MUX_LABEL                = "muxLabel"
	OPTION_KCP_LABEL                = "kcpLabel"
	OPTION_SESSION                  = "session"
	OPTION_NEW_PEER_CONNECTION_FUNC = "newPeerConnectionFunc"

	SYS_METHOD_PING                = "ping"
	SYS_METHOD_NEW_CHANNEL         = "newChannel"
	SYS_METHOD_ADD_PEER_CONNECTION = "addPeerConnection"
	SYS_METHOD_CLOSE               = "close"

	CHANNEL_MODE_RAW = "raw"
	CHANNEL_MODE_MUX = "mux"
	CHANNEL_MODE_KCP = "kcp"

	SALT = "meepo"

	IGNORE_DATA_CHANNEL_LABEL = "_ignore_"
)

func defaultNewWebrtcTransportOptions() option.Option {
	return option.NewOption(map[string]any{
		OPTION_GATHER_TIMEOUT:                61 * time.Second,
		OPTION_TEMP_DATA_CHANNEL_TIMEOUT:     17 * time.Second,
		well_known_option.OPTION_RAND_SOURCE: mrand.NewSource(time.Now().Unix()),
		transport_core.OPTION_READY_TIMEOUT:  601 * time.Second,

		well_known_option.OPTION_ENABLE_MUX:     true,
		well_known_option.OPTION_MUX_VER:        C.SMUX_VERSION,
		well_known_option.OPTION_MUX_BUF:        C.SMUX_BUFFER_SIZE,
		well_known_option.OPTION_MUX_STREAM_BUF: C.SMUX_STREAM_BUFFER_SIZE,

		well_known_option.OPTION_ENABLE_KCP:       true,
		well_known_option.OPTION_KCP_DATA_SHARD:   C.KCP_DATA_SHARD,
		well_known_option.OPTION_KCP_PARITY_SHARD: C.KCP_PARITY_SHARD,
	})
}

func defaultNewSourceTransportOptions() option.Option {
	m := defaultNewWebrtcTransportOptions()
	m[OPTION_ROLE] = "source"
	return option.NewOption(m)
}

func defaultNewSinkTransportOptions() option.Option {
	m := defaultNewWebrtcTransportOptions()
	m[OPTION_ROLE] = "sink"
	return option.NewOption(m)
}

func defaultNewChannelOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_MODE: CHANNEL_MODE_MUX,
	})
}

var (
	WithRole, GetRole                                     = option.New[string](OPTION_ROLE)
	WithTempDataChannelTimeout, GetTempDataChannelTimeout = option.New[time.Duration](OPTION_TEMP_DATA_CHANNEL_TIMEOUT)
	WithOffer, GetOffer                                   = option.New[webrtc.SessionDescription](OPTION_OFFER)
	WithAnswer, GetAnswer                                 = option.New[webrtc.SessionDescription](OPTION_ANSWER)
	WithGatherTimeout, GetGatherTimeout                   = option.New[time.Duration](OPTION_GATHER_TIMEOUT)
	WithGatherOnNewFunc, GetGatherOnNewFunc               = option.New[GatherFunc](OPTION_GATHER_ON_NEW_FUNC)
	WithGatherDoneOnNewFunc, GetGatherDoneOnNewFunc       = option.New[GatherDoneFunc](OPTION_GATHER_DONE_ON_NEW_FUNC)
	WithGatherFunc, GetGatherFunc                         = option.New[GatherFunc](OPTION_GATHER_FUNC)
	WithMuxLabel, GetMuxLabel                             = option.New[string](OPTION_MUX_LABEL)
	WithKcpLabel, GetKcpLabel                             = option.New[string](OPTION_KCP_LABEL)
	WithSession, GetSession                               = option.New[int32](OPTION_SESSION)
	WithNewPeerConnectionFunc, GetNewPeerConnectionFunc   = option.New[NewPeerConnectionFunc](OPTION_NEW_PEER_CONNECTION_FUNC)
)
