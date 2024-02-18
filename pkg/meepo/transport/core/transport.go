package transport_core

import (
	lib_registerer "github.com/PeerXu/meepo/pkg/lib/registerer"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

type Transport = meepo_interface.Transport

var RegisterNewTransportFunc, NewTransport = lib_registerer.Pair[Transport]()
