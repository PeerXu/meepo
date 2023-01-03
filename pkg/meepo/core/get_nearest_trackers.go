package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) getNearestTrackers(target Addr, count int, excludes []Addr) (tks []Tracker, found bool, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getNearestTrackers",
		"target":  target.String(),
	})

	targetID := Addr2ID(target)
	excludeIDs := Addrs2IDs(excludes)
	nearestIDs, found := mp.routingTable.NearestIDs(targetID, count, excludeIDs)

	nearestAddrs := IDs2Addrs(nearestIDs)
	nearestTrackers, err := mp.listTrackersByAddrs(nearestAddrs)
	if err != nil {
		logger.WithError(err).Debugf("failed to list trackers by addrs")
		return nil, false, err
	}

	logger.Tracef("get nearest trackers")

	return nearestTrackers, found, nil
}
