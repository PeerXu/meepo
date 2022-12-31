package meepo_core

import (
	"context"

	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIListChannels(ctx context.Context, _req any) (any, error) {
	tcs, err := mp.ListChannels(ctx)
	if err != nil {
		return nil, err
	}

	var cvs []sdk_interface.ChannelView
	for transportAddr, cs := range tcs {
		cvs = append(cvs, ViewChannelsWithAddr(cs, transportAddr)...)
	}
	return cvs, nil
}
