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
	case <-c.readyCh:
		if c.readyErr != nil {
			return c.readyErr
		}
		return nil
	case <-time.After(c.readyTimeout):
		c.readyWithError(transport_core.ErrReadyTimeout)
		return c.readyErr
	}
}
