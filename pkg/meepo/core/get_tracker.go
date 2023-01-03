package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) getTracker(target addr.Addr) (Tracker, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getTracker",
		"target":  target.String(),
	})

	mp.trackersMtx.Lock()
	defer mp.trackersMtx.Unlock()

	tk, ok := mp.trackers[target]
	if !ok {
		err := ErrTrackerNotFoundFn(target.String())
		logger.WithError(err).Debugf("tracker not found")
		return nil, err
	}

	logger.Tracef("get tracker")

	return tk, nil
}
