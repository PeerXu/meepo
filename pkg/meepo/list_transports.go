package meepo

import (
	"github.com/PeerXu/meepo/pkg/transport"
)

func (mp *Meepo) ListTransports() ([]transport.Transport, error) {
	logger := mp.getLogger().WithField("#method", "ListTransports")

	tps, err := mp.listTransports()
	if err != nil {
		logger.WithError(err).Errorf("failed to list transports")
		return nil, err
	}

	logger.Debugf("list transports")

	return tps, nil
}
