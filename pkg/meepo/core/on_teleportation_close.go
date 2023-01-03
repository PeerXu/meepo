package meepo_core

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (mp *Meepo) onTeleportationClose(tp Teleportation) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":         "onTeleportationClose",
		"teleportationID": tp.ID(),
	})

	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	delete(mp.teleportations, tp.ID())

	logger.Tracef("teleportation closed")

	return nil
}
