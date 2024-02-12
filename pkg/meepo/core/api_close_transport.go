package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiCloseTransport(ctx context.Context, req sdk_interface.CloseTransportRequest) (res rpc_core.EMPTY, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPICloseTransport",
		"target":  req.Target,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return
	}

	if err = t.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close transport")
		return
	}

	logger.Infof("transport closed")

	return rpc_core.NO_CONTENT, nil
}
