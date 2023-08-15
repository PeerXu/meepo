package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) onRemoveWebrtcTransport(t Transport) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	mp.onRemoveWebrtcTransportNL(t)
}

func (mp *Meepo) onRemoveWebrtcTransportNL(t Transport) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onRemoveWebrtcTransportNL",
		"addr":    t.Addr().String(),
	})

	mp.removeTransportNL(t.Addr())
	mp.removeTrackerNL(t.Addr())

	if tk, found := mp.defaultTrackers[t.Addr()]; found {
		mp.trackers[tk.Addr()] = tk
		logger.Tracef("addr in default trackers, replace by default tracker")
	} else {
		mp.routingTable.RemoveID(Addr2ID(t.Addr())) // nolint:errcheck
		logger.Tracef("remove addr from routing table")
	}
}

func (mp *Meepo) onRemovePipeTransport(t Transport) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	mp.onRemovePipeTransportNL(t)
}

func (mp *Meepo) onRemovePipeTransportNL(t Transport) {
	mp.removeTransportNL(t.Addr())
}
