package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/version"
)

func (mp *Meepo) hdrAPIGetVersion(ctx context.Context, _ any) (any, error) {
	return *version.Get(), nil
}
