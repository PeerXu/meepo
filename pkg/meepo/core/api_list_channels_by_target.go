package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiListChannelsByTarget(ctx context.Context, req sdk_interface.ListChannelsByTarget) (res []sdk_interface.ChannelView, err error) {
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
