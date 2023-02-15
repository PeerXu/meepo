package tracker_rpc

import (
	"github.com/pion/webrtc/v3"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func (tk *RPCTracker) AddPeerConnection(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	err = tk.caller.Call(tk.context(), tracker_core.METHOD_ADD_PEER_CONNECTION, in, &answer, well_known_option.WithDestination(tk.Addr().Bytes()))
	return
}
