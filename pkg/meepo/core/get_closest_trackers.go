package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) getClosestTrackers(target addr.Addr) ([]Tracker, bool, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getClosestTrackers",
		"target":  target.String(),
	})

	targetID := Addr2ID(target)
	closestIDs, found := mp.routingTable.ClosestIDs(targetID, mp.dhtAlpha)
	closestTrackers, err := mp.listTrackersByAddrs(IDs2Addrs(closestIDs))
	if err != nil {
		logger.WithError(err).Debugf("failed to list trackers by addrs")
		return nil, false, err
	}

	logger.Tracef("get closest trackers")

	return closestTrackers, found, nil
}
