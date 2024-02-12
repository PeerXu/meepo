package meepo_core

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiListChannels(ctx context.Context, _ rpc_core.EMPTY) (res []sdk_interface.ChannelView, err error) {
	tcs, err := mp.ListChannels(ctx)
	if err != nil {
		return nil, err
	}

	for transportAddr, cs := range tcs {
		res = append(res, ViewChannelsWithAddr(cs, transportAddr)...)
	}
	return res, nil
}
