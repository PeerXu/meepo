package meepo

import "github.com/PeerXu/meepo/pkg/transport"

func (mp *Meepo) GetTransport(peerID string) (transport.Transport, error) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.getTransportNL(peerID)
}

func (mp *Meepo) getTransportNL(peerID string) (transport.Transport, error) {
	transport, ok := mp.transports[peerID]
	if !ok {
		return nil, TransportNotExistError
	}

	return transport, nil
}

func (mp *Meepo) listTransports() ([]transport.Transport, error) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	return mp.listTransportsNL()
}

func (mp *Meepo) listTransportsNL() ([]transport.Transport, error) {
	var tps []transport.Transport

	for _, tp := range mp.transports {
		tps = append(tps, tp)
	}

	return tps, nil
}
