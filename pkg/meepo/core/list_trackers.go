package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
)

func (mp *Meepo) listTrackersByAddrs(addrs []addr.Addr) ([]Tracker, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "listTrackersByAddrs",
	})

	mp.trackersMtx.Lock()
	defer mp.trackersMtx.Unlock()

	var tks []Tracker
	for _, x := range addrs {
		tk, ok := mp.trackers[x]
		if !ok {
			err := ErrTrackerNotFoundFn(x.String())
			logger.WithError(err).Debugf("tracker not found")
			return nil, err
		}
		tks = append(tks, tk)
	}

	logger.Tracef("list trackers by addrs")

	return tks, nil
}
