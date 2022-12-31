package tracker_transport

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
)

func (tk *TransportTracker) NewTransport(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	ctx := tk.context()
	err = tk.transport.Call(ctx, "newTransport", in, &answer, well_known_option.WithScope("sys"))
	return
}
