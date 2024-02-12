package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

func (mp *Meepo) apiDiagnostic(ctx context.Context, _ rpc_core.EMPTY) (map[string]any, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPIDiagnostic",
	})

	rp, err := mp.Diagnostic(ctx)
	if err != nil {
		logger.WithError(err).Errorf("failed to diagnostic")
		return nil, err
	}

	return rp, nil
}
