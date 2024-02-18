package rpc_core

import (
	lib_registerer "github.com/PeerXu/meepo/pkg/lib/registerer"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type Server = rpc_interface.Server

var RegisterNewServerFunc, NewServer = lib_registerer.Pair[Server]()
