package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) ListTransports(ctx context.Context, opts ...ListTransportsOption) (ts []Transport, err error) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "ListTransports",
	})

	for _, t := range mp.transports {
		ts = append(ts, t)
	}

	logger.Tracef("list transports")

	return
}
