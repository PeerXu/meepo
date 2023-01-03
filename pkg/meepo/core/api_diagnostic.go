package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) hdrAPIDiagnostic(ctx context.Context, _ any) (any, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPIDiagnostic",
	})

	rp, err := mp.Diagnostic()
	if err != nil {
		logger.WithError(err).Errorf("failed to diagnostic")
		return nil, err
	}

	return rp, nil
}
