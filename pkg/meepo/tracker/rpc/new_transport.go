package tracker_rpc

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
)

func (tk *RPCTracker) NewTransport(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	ctx := tk.context()
	err = tk.caller.Call(ctx, "newTransport", in, &answer, well_known_option.WithDestination(tk.Addr().Bytes()))
	return
}
