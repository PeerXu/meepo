package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
)

func (mp *Meepo) getClosestTrackers(target addr.Addr) ([]Tracker, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getClosestTrackers",
		"target":  target.String(),
	})

	targetID := Addr2ID(target)
	closestIDs, found := mp.routingTable.ClosestIDs(targetID, mp.dhtAlpha)

	var closestAddrs []addr.Addr
	if found {
		closestAddrs = IDs2Addrs(closestIDs[:1])
	} else {
		closestAddrs = IDs2Addrs(closestIDs)
	}

	closestTrackers, err := mp.listTrackersByAddrs(closestAddrs)
	if err != nil {
		logger.WithError(err).Debugf("failed to list trackers by addrs")
		return nil, err
	}

	logger.Tracef("get closest trackers")

	return closestTrackers, nil
}
