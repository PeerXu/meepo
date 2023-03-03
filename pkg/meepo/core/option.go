package meepo_core

import (
	"time"

	C "github.com/PeerXu/meepo/pkg/lib/constant"
	"github.com/PeerXu/meepo/pkg/lib/option"
	mrand "github.com/PeerXu/meepo/pkg/lib/rand"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	meepo_routing_table_core "github.com/PeerXu/meepo/pkg/meepo/routing_table/core"
)

const (
	OPTION_DHT_ALPHA         = "dhtAlpha"
	OPTION_POOF_INTERVAL     = "poofInterval"
	OPTION_POOF_COUNT        = "poofCount"
	OPTION_ENABLE_POOF       = "enablePoof"
	OPTION_GET_TRACKERS_FUNC = "getTrackersFunc"

	METHOD_PING                = "ping"
	METHOD_PERMIT              = "permit"
	METHOD_ADD_PEER_CONNECTION = "addPeerConnection"
)

type NewMeepoOption = option.ApplyOption

type NewTransportOption = meepo_interface.NewTransportOption

type ListTransportsOption = meepo_interface.ListTransportsOption

type GetTransportOption = meepo_interface.GetTransportOption

type NewChannelOption = meepo_interface.NewChannelOption

type ListChannelsOption = meepo_interface.ListChannelsOption

type GetChannelOption = meepo_interface.GetChannelOption

type NewTeleportationOption = meepo_interface.NewTeleportationOption

type ListTeleportationsOption = meepo_interface.ListTeleportationsOption

type GetTeleportationOption = meepo_interface.GetTeleportationOption

type TeleportOption = meepo_interface.TeleportOption

var (
	WithDHTAlpha, GetDHTAlpha               = option.New[int](OPTION_DHT_ALPHA)
	WithPoofInterval, GetPoofInterval       = option.New[time.Duration](OPTION_POOF_INTERVAL)
	WithPoofCount, GetPoofCount             = option.New[int](OPTION_POOF_COUNT)
	WithEnablePoof, GetEnablePoof           = option.New[bool](OPTION_ENABLE_POOF)
	WithGetTrackersFunc, GetGetTrackersFunc = option.New[getTrackersFunc](OPTION_GET_TRACKERS_FUNC)
)

func defaultNewMeepoOptions() option.Option {
	return option.NewOption(map[string]any{
		OPTION_ENABLE_POOF:                         true,
		OPTION_DHT_ALPHA:                           8,
		OPTION_POOF_INTERVAL:                       31 * time.Second,
		OPTION_POOF_COUNT:                          3,
		meepo_routing_table_core.OPTION_GREEN_LINE: 5,
		well_known_option.OPTION_RAND_SOURCE:       mrand.NewSource(time.Now().UnixNano()),

		well_known_option.OPTION_WEBRTC_RECEIVE_BUFFER_SIZE: uint32(33554432),

		well_known_option.OPTION_ENABLE_MUX:     true,
		well_known_option.OPTION_MUX_VER:        C.SMUX_VERSION,
		well_known_option.OPTION_MUX_BUF:        C.SMUX_BUFFER_SIZE,
		well_known_option.OPTION_MUX_STREAM_BUF: C.SMUX_STREAM_BUFFER_SIZE,
		well_known_option.OPTION_MUX_NOCOMP:     C.SMUX_NOCOMP,

		well_known_option.OPTION_ENABLE_KCP:       true,
		well_known_option.OPTION_KCP_PRESET:       C.KCP_PRESET,
		well_known_option.OPTION_KCP_CRYPT:        C.KCP_CRYPT,
		well_known_option.OPTION_KCP_KEY:          C.KCP_KEY,
		well_known_option.OPTION_KCP_MTU:          C.KCP_MTU,
		well_known_option.OPTION_KCP_SNDWND:       C.KCP_SNDWND,
		well_known_option.OPTION_KCP_RCVWND:       C.KCP_RCVWND,
		well_known_option.OPTION_KCP_DATA_SHARD:   C.KCP_DATA_SHARD,
		well_known_option.OPTION_KCP_PARITY_SHARD: C.KCP_PARITY_SHARD,
	})
}

func (mp *Meepo) defaultNewTransportOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_ENABLE_MUX:     mp.enableMux,
		well_known_option.OPTION_MUX_VER:        mp.muxVer,
		well_known_option.OPTION_MUX_BUF:        mp.muxBuf,
		well_known_option.OPTION_MUX_STREAM_BUF: mp.muxStreamBuf,
		well_known_option.OPTION_MUX_NOCOMP:     mp.muxNocomp,

		well_known_option.OPTION_ENABLE_KCP:       mp.enableKcp,
		well_known_option.OPTION_KCP_PRESET:       mp.kcpPreset,
		well_known_option.OPTION_KCP_CRYPT:        mp.kcpCrypt,
		well_known_option.OPTION_KCP_KEY:          mp.kcpKey,
		well_known_option.OPTION_KCP_MTU:          mp.kcpMtu,
		well_known_option.OPTION_KCP_SNDWND:       mp.kcpSndwnd,
		well_known_option.OPTION_KCP_RCVWND:       mp.kcpRcvwnd,
		well_known_option.OPTION_KCP_DATA_SHARD:   mp.kcpDataShard,
		well_known_option.OPTION_KCP_PARITY_SHARD: mp.kcpParityShard,
	})
}
