package meepo

import "github.com/PeerXu/meepo/pkg/teleportation"

func (mp *Meepo) ListTeleportations() ([]teleportation.Teleportation, error) {
	logger := mp.getLogger().WithField("#method", "ListTeleportations")

	tss, err := mp.listTeleportations()
	if err != nil {
		logger.WithError(err).Errorf("failed to list teleportations")
		return nil, err
	}

	logger.Debugf("list teleportations")

	return tss, nil
}
