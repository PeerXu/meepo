package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiGetTransport(ctx context.Context, req sdk_interface.GetTransportRequest) (res sdk_interface.TransportView, err error) {
	target, err := addr.FromString(req.Target)
	if err != nil {
		return
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		return
	}

	return ViewTransport(t), nil
}
