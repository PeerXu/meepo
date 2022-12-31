package transport_webrtc

import (
	"github.com/PeerXu/meepo/pkg/internal/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) Handle(method string, fn meepo_interface.HandleFunc, opts ...meepo_interface.HandleOption) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#mwethod": "Handle",
		"method":   method,
	})

	t.fnsMtx.Lock()
	defer t.fnsMtx.Unlock()

	t.fns[method] = fn

	logger.Tracef("handle method")
}
