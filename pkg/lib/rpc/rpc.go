package rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	_ "github.com/PeerXu/meepo/pkg/lib/rpc/default"
	_ "github.com/PeerXu/meepo/pkg/lib/rpc/http"
)

var (
	NewCaller  = rpc_core.NewCaller
	NewServer  = rpc_core.NewServer
	NewHandler = rpc_core.NewHandler
)
