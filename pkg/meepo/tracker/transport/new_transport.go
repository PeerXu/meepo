package tracker_transport

import (
	"github.com/pion/webrtc/v3"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func (tk *TransportTracker) NewTransport(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	err = tk.transport.Call(tk.context(), tracker_core.METHOD_NEW_TRANSPORT, in, &answer)
	return
}
