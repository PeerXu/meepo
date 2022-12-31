package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) GetTransport(ctx context.Context, target Addr, opts ...GetTransportOption) (Transport, error) {
	var err error

	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "GetTransport",
		"target":  target,
	})

	t, ok := mp.transports[target]
	if !ok {
		err = ErrTransportNotFoundFn(target.String())
		logger.WithError(err).Debugf("transport not found")
		return nil, err
	}

	logger.Tracef("get transport")

	return t, nil
}
