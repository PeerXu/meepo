package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIGetChannel(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.GetChannelRequest)

	target, err := addr.FromString(req.Target)
	if err != nil {
		return nil, err
	}

	c, err := mp.GetChannel(ctx, target, req.ChannelID)
	if err != nil {
		return nil, err
	}

	return ViewChannelWithAddr(c, target), nil
}
