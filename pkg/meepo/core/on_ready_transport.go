package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/meepo/tracker"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (mp *Meepo) onReadyWebrtcTransport(t Transport) error {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.onReadyWebrtcTransportNL(t)
}

func (mp *Meepo) onReadyWebrtcTransportNL(t Transport) (err error) {
	if err = mp.routingTable.AddID(Addr2ID(t.Addr())); err != nil {
		return err
	}

	tk, err := tracker.NewTracker("transport",
		transport_core.WithTransport(t),
	)
	if err != nil {
		return err
	}

	mp.trackers[tk.Addr()] = tk

	return nil
}
