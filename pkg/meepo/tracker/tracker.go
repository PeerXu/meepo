package tracker

import (
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	_ "github.com/PeerXu/meepo/pkg/meepo/tracker/rpc"
	_ "github.com/PeerXu/meepo/pkg/meepo/tracker/transport"
)

var NewTracker = tracker_core.NewTracker
