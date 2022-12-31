package transport

import (
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
	_ "github.com/PeerXu/meepo/pkg/meepo/transport/pipe"
	_ "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

var (
	NewTransport = transport_core.NewTransport
)
