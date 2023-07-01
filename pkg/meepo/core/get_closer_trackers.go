package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) getCloserTrackers(target Addr, requestCandidates int, excludes []Addr) (tks []Tracker, found bool, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getCloserTrackers",
		"target":  target.String(),
	})

	targetID := Addr2ID(target)
	excludeIDs := Addrs2IDs(excludes)
	closerIDs, found := mp.routingTable.CloserIDs(targetID, requestCandidates, excludeIDs)

	closerAddrs := IDs2Addrs(closerIDs)
	closerTrackers, err := mp.listTrackersByAddrs(closerAddrs)
	if err != nil {
		logger.WithError(err).Debugf("failed to list trackers by addrs")
		return nil, false, err
	}

	logger.Tracef("get closer trackers")

	return closerTrackers, found, nil
}
