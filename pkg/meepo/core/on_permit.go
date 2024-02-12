package meepo_core

import (
	"context"
	"net"

	"github.com/PeerXu/meepo/pkg/lib/acl"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type PermitRequest struct {
	Network string
	Address string
}

func (mp *Meepo) onPermit(ctx context.Context, req PermitRequest) (res rpc_core.EMPTY, err error) {
	t := transport_core.ContextGetTransport(ctx)

	err = mp.permit(t.Addr().String(), req.Network, req.Address)
	if err != nil {
		return
	}

	return
}

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
