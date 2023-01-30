package meepo_core

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (mp *Meepo) onRemoveWebrtcTransport(t Transport) error {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.onRemoveWebrtcTransportNL(t)
}

func (mp *Meepo) onRemoveWebrtcTransportNL(t Transport) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onRemoveWebrtcTransportNL",
		"addr":    t.Addr().String(),
	})

	delete(mp.transports, t.Addr())
	logger.Tracef("delete transport")

	delete(mp.trackers, t.Addr())
	logger.Tracef("delete tracker")

	if tk, found := mp.defaultTrackers[t.Addr()]; found {
		mp.trackers[tk.Addr()] = tk
		logger.Tracef("addr in default trackers, replace by default tracker")
	} else {
		mp.routingTable.RemoveID(Addr2ID(t.Addr())) // nolint:errcheck
		logger.Tracef("remove addr from routing table")
	}

	return nil
}

func (mp *Meepo) onRemovePipeTransport(t Transport) error {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.onRemovePipeTransportNL(t)
}

func (mp *Meepo) onRemovePipeTransportNL(t Transport) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onRemovePipeTransportNL",
		"addr":    t.Addr().String(),
	})

	delete(mp.transports, t.Addr())
	logger.Tracef("delete transport")

	return nil
}
