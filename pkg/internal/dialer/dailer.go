package dialer

import (
	dialer_core "github.com/PeerXu/meepo/pkg/internal/dialer/core"
	_ "github.com/PeerXu/meepo/pkg/internal/dialer/net"
)

type Dialer = dialer_core.Dialer

var GetGlobalDialer = dialer_core.GetGlobalDialer
