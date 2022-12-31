package transport_webrtc

import (
	"time"

	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (c *WebrtcChannel) WaitReady() error {
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
			c.readyErr = transport_core.ErrReadyTimeout
			close(c.ready)
		})
		return c.readyErr
	}
}
