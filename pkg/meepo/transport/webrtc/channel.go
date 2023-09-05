package transport_webrtc

import (
	"net"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"

	matomic "github.com/PeerXu/meepo/pkg/lib/atomic"
	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type WebrtcChannel struct {
	id       uint16
	sinkAddr net.Addr
	logger   logging.Logger

	dc *webrtc.DataChannel
	s  matomic.GenericValue[meepo_interface.ChannelState]

	readyErr     error
	readyTimeout time.Duration
	readyCh      chan struct{}
	readyOnce    sync.Once

	mode string

	onStateChange          transport_core.OnChannelStateChangeFunc
	beforeCloseChannelHook transport_core.BeforeCloseChannelHook
	afterCloseChannelHook  transport_core.AfterCloseChannelHook
}

type WebrtcSourceChannel struct {
	*WebrtcChannel
	conn meepo_interface.Conn
}

type WebrtcSinkChannel struct {
	*WebrtcChannel
	upstream      dialer_interface.Conn
	downstreamVal matomic.GenericValue[meepo_interface.Conn]
}
