package well_known_option

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_PEER_CONNECTION = "peerConnection"
)

var (
	WithPeerConnection, GetPeerConnection = option.New[*webrtc.PeerConnection](OPTION_PEER_CONNECTION)
)
