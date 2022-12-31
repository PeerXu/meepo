package transport_webrtc

import (
	"time"

	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) WaitReady() error {
	if t.readyErr != nil {
		return t.readyErr
	}

	select {
	case <-t.ready:
		if t.readyErr != nil {
			return t.readyErr
		}
		return nil
	case <-time.After(t.readyTimeout):
		t.readyOnce.Do(func() {
			close(t.ready)
			t.readyErr = transport_core.ErrReadyTimeout
		})
		return t.readyErr
	}
}
