package meepo_core

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (mp *Meepo) removeTrackerNL(addr Addr) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "removeTrackerNL",
		"addr":    addr.String(),
	})
	delete(mp.trackers, addr)
	logger.Tracef("remove tracker")
}
