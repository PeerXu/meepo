package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type (
	PingRequest struct {
		Nonce uint32
	}

	PingResponse struct {
		Nonce uint32
	}
)

func (mp *Meepo) onPing(ctx context.Context, req PingRequest) (res PingResponse, err error) {
	t := transport_core.ContextGetTransport(ctx)
	from := t.Addr()
	nonce := req.Nonce
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onPing",
		"from":    from.String(),
		"nonce":   nonce,
	})

	logger.Tracef("ping")
	return PingResponse{
		Nonce: nonce,
	}, nil
}
