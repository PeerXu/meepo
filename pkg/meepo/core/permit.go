package meepo_core

import (
	"net"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/acl"
)

func (mp *Meepo) permit(target, network, address string) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "permit",
	})
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		logger.WithError(err).Debugf("failed to split host port")
		return err
	}
	challenge := acl.NewEntity(target, network, host, port)
	logger = logger.WithField("challenge", challenge.String())

	if err = mp.acl.Permit(challenge); err != nil {
		logger.WithError(err).Debugf("not permitted")
		return err
	}

	logger.Tracef("permitted")

	return nil
}
