package transport_pipe

import (
	"time"

	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (c *PipeChannel) WaitReady() error {
	if c.readyErr != nil {
		return c.readyErr
	}

	select {
	case <-c.ready:
		if c.readyErr != nil {
			return c.readyErr
		}
		return nil
	case <-time.After(c.readyTimeout):
		c.readyOnce.Do(func() {
			close(c.ready)
			c.readyErr = transport_core.ErrReadyTimeout
		})
		return c.readyErr
	}
}
