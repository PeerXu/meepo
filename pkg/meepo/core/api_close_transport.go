package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPICloseTransport(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.CloseTransportRequest)

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPICloseTransport",
		"target":  req.Target,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return nil, err
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return nil, err
	}

	if err = t.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close transport")
		return nil, err
	}

	logger.Infof("transport closed")

	return rpc_core.NO_CONTENT(), nil
}
