package config

import (
	C "github.com/PeerXu/meepo/pkg/lib/constant"
	"github.com/PeerXu/meepo/pkg/lib/stun"
)

type Config struct {
	Meepo Meepo `yaml:"meepo"`
}

func (c *Config) Init() {
	if len(c.Meepo.Trackerds) == 0 && c.Meepo.Trackerd != nil {
		c.Meepo.Trackerds = append(c.Meepo.Trackerds, c.Meepo.Trackerd)
	}

	if len(c.Meepo.Trackers) == 0 && c.Meepo.Tracker != nil {
		c.Meepo.Trackers = append(c.Meepo.Trackers, c.Meepo.Tracker)
	}
}

var cfg Config

func Get() *Config {
	return &cfg
}

func Default() *Config {
	return &Config{
		Meepo: Meepo{
			Daemon: true,
			Pprof:  "",
			Identity: Identity{
				NoFile: false,
				File:   "",
			},
			Tracker: &Tracker{
				Name:       "rpc",
				CallerName: "http",
				Addr:       C.TRACKER_ADDR,
				Host:       C.TRACKER_HOST,
			},
			API: API{
				Name: "http",
				Host: C.API_HOST,
			},
			Socks5: Socks5{
				Host: C.SOCKS5_HOST,
			},
			Acl: C.ACL,
			Log: Log{
				Level: C.LOG_LEVEL,
			},
			Webrtc: Webrtc{
				IceServers:     stun.STUNS,
				RecvBufferSize: C.WEBRTC_RECEIVE_BUFFER_SIZE,
			},
			Smux: Smux{
				Disable:          false,
				Version:          C.SMUX_VERSION,
				BufferSize:       C.SMUX_BUFFER_SIZE,
				StreamBufferSize: C.SMUX_STREAM_BUFFER_SIZE,
				Nocomp:           C.SMUX_NOCOMP,
			},
			Kcp: Kcp{
				Disable:     true,
				Preset:      C.KCP_PRESET,
				Crypt:       C.KCP_CRYPT,
				Key:         C.KCP_KEY,
				Mtu:         C.KCP_MTU,
				Sndwnd:      C.KCP_SNDWND,
				Rcvwnd:      C.KCP_RCVWND,
				DataShard:   C.KCP_DATA_SHARD,
				ParityShard: C.KCP_PARITY_SHARD,
			},
		},
	}
}
