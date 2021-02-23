package meepo

import "github.com/sirupsen/logrus"

func (mp *Meepo) CloseTransport(peerID string) error {
	return mp.closeTransport(peerID)
}

func (mp *Meepo) closeTransport(peerID string) error {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "closeTransport",
		"peerID":  peerID,
	})

	tp, err := mp.GetTransport(peerID)
	if err != nil {
		logger.WithError(err).Errorf("transport not found")
		return err
	}

	if err = tp.Close(); err != nil {
		logger.WithError(err).Errorf("failed to close transport")
		return err
	}

	mp.removeTransport(peerID)
	logger.Tracef("remove transport")
	mp.removeTeleportationsByPeerID(peerID)
	logger.Tracef("remove teleportations")

	logger.Infof("transport closed")

	return nil
}
