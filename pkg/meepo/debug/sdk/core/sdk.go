package meepo_debug_sdk_core

import (
	lib_registerer "github.com/PeerXu/meepo/pkg/lib/registerer"
	meepo_debug_interface "github.com/PeerXu/meepo/pkg/meepo/debug/interface"
)

var Register, New = lib_registerer.Pair[meepo_debug_interface.MeepoDebugInterface]()
