package listenerer_http

import (
	"time"

	listenerer_core "github.com/PeerXu/meepo/pkg/lib/listenerer/core"
)

func (c *HttpConn) WaitEnabled(timeout time.Duration) error {
	select {
	case <-c.enable:
		return nil
	case <-c.close:
		return listenerer_core.ErrConnClosed
	case <-time.After(timeout):
		return listenerer_core.ErrConnWaitEnabledTimeout
	}
}
