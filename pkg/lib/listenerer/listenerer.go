package listenerer

import (
	listenerer_core "github.com/PeerXu/meepo/pkg/lib/listenerer/core"
	_ "github.com/PeerXu/meepo/pkg/lib/listenerer/net"
	_ "github.com/PeerXu/meepo/pkg/lib/listenerer/socks5"
)

type Listenerer = listenerer_core.Listenerer

var GetGlobalListenerer = listenerer_core.GetGlobalListenerer
