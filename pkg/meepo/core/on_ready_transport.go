package meepo_core

import (
	"errors"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/routing_table"
	"github.com/PeerXu/meepo/pkg/meepo/tracker"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (mp *Meepo) onReadyWebrtcTransport(t Transport) {
	mp.trackersMtx.Lock()
	defer mp.trackersMtx.Unlock()

	mp.onReadyWebrtcTransportNL(t)
}

func (mp *Meepo) onReadyWebrtcTransportNL(t Transport) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onReadyWebrtcTransportNL",
		"addr":    t.Addr().String(),
	})

	if err := mp.routingTable.AddID(Addr2ID(t.Addr())); err != nil {
		if !errors.Is(err, routing_table.ErrOutOfBucketSize) {
			logger.WithError(err).Debugf("failed to add addr to routing table")
			return
		}

		logger.Debugf("routing table is full, drop tracker")
	}

	tk, _ := tracker.NewTracker("transport", transport_core.WithTransport(t))
	mp.trackers[tk.Addr()] = tk
}
