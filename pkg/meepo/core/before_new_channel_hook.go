package meepo_core

import (
	"github.com/PeerXu/meepo/pkg/lib/option"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (mp *Meepo) beforeNewChannelHook(t Transport, network, address string, opts ...transport_core.HookOption) error {
	o := option.Apply(opts...)

	isSink, _ := transport_core.GetIsSink(o)
	if isSink {
		return mp.permit(t.Addr().String(), network, address)
	}

	return nil
}
