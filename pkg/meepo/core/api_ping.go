package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIPing(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.PingRequest)

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPIPing",
		"target":  req.Target,
		"nonce":   req.Nonce,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target addr")
		return nil, err
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return nil, err
	}

	var res PingResponse
	if err = t.Call(ctx, "ping", &PingRequest{Nonce: req.Nonce}, &res); err != nil {
		logger.WithError(err).Errorf("failed to ping")
		return nil, err
	}

	if res.Nonce != req.Nonce {
		err = ErrInvalidNonceFn(res.Nonce)
		logger.WithError(err).Errorf("invalid nonce")
		return nil, err
	}

	return &res, nil
}
