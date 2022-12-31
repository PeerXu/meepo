package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) GetTeleportation(ctx context.Context, id string, opts ...GetTeleportationOption) (Teleportation, error) {
	var err error

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "GetTeleportation",
		"id":      id,
	})

	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	tp, ok := mp.teleportations[id]
	if !ok {
		err = ErrTeleportationNotFoundFn(id)
		logger.WithError(err).Debugf("teleportation not found")
		return nil, err
	}

	logger.Tracef("get teleportation")

	return tp, nil
}
