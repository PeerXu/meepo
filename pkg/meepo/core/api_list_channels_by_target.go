package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIListChannelsByTarget(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.ListChannelsByTarget)
	target, err := addr.FromString(req.Target)
	if err != nil {
		return nil, err
	}

	cs, err := mp.ListChannelsByTarget(ctx, target)
	if err != nil {
		return nil, err
	}

	return ViewChannelsWithAddr(cs, target), nil
}
