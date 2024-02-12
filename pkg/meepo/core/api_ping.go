package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiPing(ctx context.Context, req sdk_interface.PingRequest) (res sdk_interface.PingResponse, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "apiPing",
		"target":  req.Target,
		"nonce":   req.Nonce,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target addr")
		return
	}

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return
	}

	if err = t.Call(ctx, METHOD_PING, &PingRequest{Nonce: req.Nonce}, &res); err != nil {
		logger.WithError(err).Errorf("failed to ping")
		return
	}

	if res.Nonce != req.Nonce {
		err = ErrInvalidNonceFn(res.Nonce)
		logger.WithError(err).Errorf("invalid nonce")
		return
	}

	return
}
