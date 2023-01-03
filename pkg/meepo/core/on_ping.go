package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
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

func (mp *Meepo) newOnPingRequest() any { return &PingRequest{} }

func (mp *Meepo) hdrOnPing(ctx context.Context, _req any) (any, error) {
	req := _req.(*PingRequest)
	t := transport_core.ContextGetTransport(ctx)

	nonce, err := mp.onPing(t.Addr(), req.Nonce)
	if err != nil {
		return nil, err
	}

	res := &PingResponse{Nonce: nonce}

	return res, nil
}

func (mp *Meepo) onPing(from addr.Addr, nonce uint32) (uint32, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onPing",
		"from":    from.String(),
		"nonce":   nonce,
	})

	logger.Tracef("ping")

	return nonce, nil
}
