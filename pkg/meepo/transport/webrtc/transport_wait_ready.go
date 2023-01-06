package transport_webrtc

import (
	"time"

	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) WaitReady() error {
	if err := t.readyError(); err != nil {
		return err
	}

	select {
	case <-t.ready:
		if err := t.readyError(); err != nil {
			return err
		}

		return nil
	case <-time.After(t.readyTimeout):
		t.readyOnce.Do(func() {
			close(t.ready)
			t.readyErrVal.Store(transport_core.ErrReadyTimeout)
		})
		return t.readyError()
	}
}
