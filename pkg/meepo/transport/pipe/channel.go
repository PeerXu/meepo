package transport_pipe

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type PipeChannel struct {
	id       uint16
	state    atomic.Value
	conn     meepo_interface.Conn
	sinkAddr net.Addr
	logger   logging.Logger

	beforeCloseChannelHook transport_core.BeforeCloseChannelHook
	afterCloseChannelHook  transport_core.AfterCloseChannelHook

	readyErr     error
	readyTimeout time.Duration
	ready        chan struct{}
	readyOnce    sync.Once
}
