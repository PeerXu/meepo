package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/dialer"
	"github.com/PeerXu/meepo/pkg/internal/listenerer"
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	teleportation_core "github.com/PeerXu/meepo/pkg/meepo/teleportation/core"
)

func (mp *Meepo) NewTeleportation(ctx context.Context, target Addr, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress string, opts ...NewTeleportationOption) (Teleportation, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":       "NewTeleportation",
		"target":        target.String(),
		"sourceNetwork": sourceNetwork,
		"sourceAddress": sourceAddress,
		"sinkNetwork":   sinkNetwork,
		"sinkAddress":   sinkAddress,
	})

	o := option.Apply(opts...)
	mode, err := well_known_option.GetMode(o)
	if err != nil {
		return nil, err
	}
	if mp.Addr().Equal(target) {
		mode = "raw"
	}
	logger = logger.WithField("mode", mode)

	id := mp.newTeleportationID()
	sinkAddr := dialer.NewAddr(sinkNetwork, sinkAddress)

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return nil, err
	}

	if err = t.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait transport ready")
		return nil, err
	}

	lis, err := listenerer.GetGlobalListenerer().Listen(
		ctx, sourceNetwork, sourceAddress,
		well_known_option.WithLogger(mp.GetRawLogger()),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to listen")
		return nil, err
	}
	sourceAddr := lis.Addr()

	tp, err := teleportation_core.NewTeleportation(
		well_known_option.WithID(id),
		well_known_option.WithMode(mode),
		well_known_option.WithLogger(mp.GetRawLogger()),
		listenerer.WithListener(lis),
		well_known_option.WithAddr(target),
		teleportation_core.WithSourceAddr(sourceAddr),
		teleportation_core.WithSinkAddr(sinkAddr),
		teleportation_core.WithOnTeleportationAcceptFunc(mp.onTeleportationAccept),
		teleportation_core.WithOnTeleportationCloseFunc(mp.onTeleportationClose),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to new teleportation")
		return nil, err
	}

	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.teleportations[tp.ID()] = tp

	logger.Tracef("new teleportation")

	return tp, nil
}
