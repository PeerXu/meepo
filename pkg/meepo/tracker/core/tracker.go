package tracker_core

import (
	lib_registerer "github.com/PeerXu/meepo/pkg/lib/registerer"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

type Tracker = tracker_interface.Tracker

var RegisterNewTrackerFunc, NewTracker = lib_registerer.Pair[Tracker]()
