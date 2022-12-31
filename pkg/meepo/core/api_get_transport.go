package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIGetTransport(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.GetTransportRequest)

	target, err := addr.FromString(req.Target)
	if err != nil {
		return nil, err
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		return nil, err
	}

	return ViewTransport(t), nil
}
