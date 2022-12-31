package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) ListTeleportations(ctx context.Context, opts ...ListTeleportationsOption) (tps []Teleportation, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "ListTeleportations",
	})

	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	for _, tp := range mp.teleportations {
		tps = append(tps, tp)
	}

	logger.Tracef("list teleportations")

	return
}
