package meepo

import (
	"errors"
	"net"
	"sync"

	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/transport"
)

func newTeleportOptions() objx.Map {
	return objx.New(map[string]interface{}{})
}

func (mp *Meepo) Teleport(peerID string, remote net.Addr, opts ...TeleportOption) (net.Addr, error) {
	var local net.Addr
	var name string
	var tp transport.Transport
	var err error
	var ok bool

	o := newTeleportOptions()

	for _, opt := range opts {
		opt(o)
	}

	tp, err = mp.GetTransport(peerID)
	if err != nil {
		if !errors.Is(err, TransportNotExistError) {
			return nil, err
		}

		var wg sync.WaitGroup
		tp, err = mp.NewTransport(peerID)
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		tp.OnTransportState(transport.TransportStateConnected, func(hid int64) {
			wg.Done()
			tp.UnsetOnTransportState(transport.TransportStateConnected, hid)
		})

		wg.Wait()
	}

	tss, err := mp.listTeleportationsByPeerID(peerID)
	if err != nil {
		return nil, err
	}

	for _, ts := range tss {
		sink := ts.Sink()
		if sink.Network() == remote.Network() &&
			sink.String() == remote.String() {
			return ts.Source(), nil
		}
	}

	var ntOpts []NewTeleportationOption
	if local, ok = o.Get("local").Inter().(net.Addr); ok {
		ntOpts = append(ntOpts, WithLocalAddress(local))
	}

	if name, ok = o.Get("name").Inter().(string); ok {
		ntOpts = append(ntOpts, WithName(name))
	}

	ts, err := mp.NewTeleportation(peerID, remote, ntOpts...)
	if err != nil {
		return nil, err
	}

	return ts.Source(), nil
}
