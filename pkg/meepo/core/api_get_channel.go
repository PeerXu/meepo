package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiGetChannel(ctx context.Context, req sdk_interface.GetChannelRequest) (res sdk_interface.ChannelView, err error) {
	target, err := addr.FromString(req.Target)
	if err != nil {
		return
	}

	c, err := mp.GetChannel(ctx, target, req.ChannelID)
	if err != nil {
		return
	}

	return ViewChannelWithAddr(c, target), nil
}
