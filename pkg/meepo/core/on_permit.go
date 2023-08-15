package meepo_core

import (
	"context"
	"net"

	"github.com/PeerXu/meepo/pkg/lib/acl"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type PermitRequest struct {
	Network string
	Address string
}

func (mp *Meepo) newOnPermitRequest() any { return &PermitRequest{} }

func (mp *Meepo) hdrOnPermit(ctx context.Context, _req any) (any, error) {
	req := _req.(*PermitRequest)
	t := transport_core.ContextGetTransport(ctx)

	err := mp.onPermit(t.Addr().String(), req.Network, req.Address)
	if err != nil {
		return nil, err
	}

	return rpc_core.NO_CONTENT(), nil
}

func (mp *Meepo) onPermit(target, network, address string) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onPermit",
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

func (mp *Meepo) beforeNewChannelHook(t Transport, network, address string, opts ...transport_core.HookOption) error {
	o := option.Apply(opts...)

	isSink, _ := transport_core.GetIsSink(o)
	if isSink {
		return mp.onPermit(t.Addr().String(), network, address)
	}

	return nil
}
