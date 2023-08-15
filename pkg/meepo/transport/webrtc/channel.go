package transport_webrtc

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pion/webrtc/v3"

	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type WebrtcChannel struct {
	id       uint16
	sinkAddr net.Addr
	logger   logging.Logger

	readyErr     error
	readyTimeout time.Duration
	ready        chan struct{}
	readyOnce    sync.Once

	mode string
}

type WebrtcSourceChannel struct {
	*WebrtcChannel
	dc   *webrtc.DataChannel
	rwc  dialer_interface.Conn
	conn meepo_interface.Conn

	beforeCloseChannelHook transport_core.BeforeCloseChannelHook
	afterCloseChannelHook  transport_core.AfterCloseChannelHook
}

type WebrtcSinkChannel struct {
	*WebrtcChannel
	dc      *webrtc.DataChannel
	rwc     dialer_interface.Conn
	connVal atomic.Value

	beforeCloseChannelHook transport_core.BeforeCloseChannelHook
	afterCloseChannelHook  transport_core.AfterCloseChannelHook
}
