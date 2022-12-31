package tracker_interface

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
)

type Tracker interface {
	Addr() addr.Addr
	NewTransport(*crypto_core.Packet) (webrtc.SessionDescription, error)
	GetCandidates(target addr.Addr, count int, excludes []addr.Addr) ([]addr.Addr, error)
}
