package transport_webrtc

import (
	"time"

	"github.com/pion/webrtc/v3"

	C "github.com/PeerXu/meepo/pkg/internal/constant"
	"github.com/PeerXu/meepo/pkg/internal/option"
	mrand "github.com/PeerXu/meepo/pkg/internal/rand"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type GatherFunc func(webrtc.SessionDescription) (webrtc.SessionDescription, error)
type GatherDoneFunc func(webrtc.SessionDescription, error)

const (
	OPTION_TEMP_DATA_CHANNEL_TIMEOUT = "tempDataChannelTimeout"
	OPTION_OFFER                     = "offer"
	OPTION_ANSWER                    = "answer"
	OPTION_GATHER_TIMEOUT            = "gatherTimeout"
	OPTION_GATHER_FUNC               = "gatherFunc"
	OPTION_GATHER_DONE_FUNC          = "gatherDoneFunc"
	OPTION_MUX_LABEL                 = "muxLabel"
	OPTION_KCP_LABEL                 = "kcpLabel"

	CHANNEL_MODE_RAW = "raw"
	CHANNEL_MODE_MUX = "mux"
	CHANNEL_MODE_KCP = "kcp"

	SALT = "meepo"
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
	return option.NewOption(m)
}

func defaultNewSinkTransportOptions() option.Option {
	m := defaultNewWebrtcTransportOptions()
	return option.NewOption(m)
}

func defaultNewChannelOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_MODE: CHANNEL_MODE_MUX,
	})
}

var (
	WithTempDataChannelTimeout, GetTempDataChannelTimeout = option.New[time.Duration](OPTION_TEMP_DATA_CHANNEL_TIMEOUT)
	WithOffer, GetOffer                                   = option.New[webrtc.SessionDescription](OPTION_OFFER)
	WithAnswer, GetAnswer                                 = option.New[webrtc.SessionDescription](OPTION_ANSWER)
	WithGatherTimeout, GetGatherTimeout                   = option.New[time.Duration](OPTION_GATHER_TIMEOUT)
	WithGatherFunc, GetGatherFunc                         = option.New[GatherFunc](OPTION_GATHER_FUNC)
	WithGatherDoneFunc, GetGatherDoneFunc                 = option.New[GatherDoneFunc](OPTION_GATHER_DONE_FUNC)
	WithMuxLabel, GetMuxLabel                             = option.New[string](OPTION_MUX_LABEL)
	WithKcpLabel, GetKcpLabel                             = option.New[string](OPTION_KCP_LABEL)
)
