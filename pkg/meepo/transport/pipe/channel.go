package transport_pipe

import (
	"net"
	"sync"
	"time"

	matomic "github.com/PeerXu/meepo/pkg/lib/atomic"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type PipeChannel struct {
	id       uint16
	state    matomic.GenericValue[meepo_interface.ChannelState]
	conn     meepo_interface.Conn
	sinkAddr net.Addr
	logger   logging.Logger

	onStateChange          transport_core.OnChannelStateChangeFunc
	beforeCloseChannelHook transport_core.BeforeCloseChannelHook
	afterCloseChannelHook  transport_core.AfterCloseChannelHook

	readyErr     error
	readyTimeout time.Duration
	ready        chan struct{}
	readyOnce    sync.Once
}
