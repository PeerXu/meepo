package meepo_core

func (mp *Meepo) onRemoveWebrtcTransport(t Transport) error {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.onRemoveWebrtcTransportNL(t)
}

func (mp *Meepo) onRemoveWebrtcTransportNL(t Transport) error {
	delete(mp.transports, t.Addr())
	delete(mp.trackers, t.Addr())
	if tk, found := mp.defaultTrackers[t.Addr()]; found {
		mp.trackers[tk.Addr()] = tk
	} else {
		mp.routingTable.RemoveID(Addr2ID(t.Addr())) // nolint:errcheck
	}

	return nil
}

func (mp *Meepo) onRemovePipeTransport(t Transport) error {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.onRemovePipeTransportNL(t)
}

func (mp *Meepo) onRemovePipeTransportNL(t Transport) error {
	delete(mp.transports, t.Addr())
	return nil
}
