package meepo_core

import (
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (mp *Meepo) onAddWebrtcTransportNL(t Transport) (err error) {
	defer func() {
		if err != nil {
			mp.onRemoveWebrtcTransportNL(t) // nolint:errcheck
		}
	}()

	t.Handle("sys/newTransport", transport_core.WrapHandleFunc(mp.newOnNewTransportRequest, mp.hdrOnNewTransport))
	t.Handle("sys/getCandidates", transport_core.WrapHandleFunc(mp.newOnGetCandidatesRequest, mp.hdrOnGetCandidates))

	t.Handle("ping", transport_core.WrapHandleFunc(mp.newOnPingRequest, mp.hdrOnPing))
	t.Handle("permit", transport_core.WrapHandleFunc(mp.newOnPermitRequest, mp.hdrOnPermit))

	mp.transports[t.Addr()] = t

	return nil
}

func (mp *Meepo) onAddPipeTransportNL(t Transport) error {
	t.Handle("ping", transport_core.WrapHandleFunc(mp.newOnPingRequest, mp.hdrOnPing))
	t.Handle("permit", transport_core.WrapHandleFunc(mp.newOnPermitRequest, mp.hdrOnPermit))

	mp.transports[t.Addr()] = t
	return nil
}
