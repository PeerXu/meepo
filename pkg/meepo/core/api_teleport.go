package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiTeleport(ctx context.Context, req sdk_interface.TeleportRequest) (res sdk_interface.TeleportationView, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":       "hdrAPITeleport",
		"target":        req.Target,
		"mode":          req.Mode,
		"sourceNetwork": req.SourceNetwork,
		"sourceAddress": req.SourceAddress,
		"sinkNetwork":   req.SinkNetwork,
		"sinkAddress":   req.SinkAddress,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return
	}

	opts := []TeleportOption{
		well_known_option.WithMode(req.Mode),
	}
	if req.Manual {
		opts = []TeleportOption{
			well_known_option.WithEnableMux(req.EnableMux),
			well_known_option.WithEnableKcp(req.EnableKcp),
		}

		if req.EnableMux {
			opts = append(opts,
				well_known_option.WithMuxVer(req.MuxVer),
				well_known_option.WithMuxBuf(req.MuxBuf),
				well_known_option.WithMuxStreamBuf(req.MuxStreamBuf),
				well_known_option.WithMuxNocomp(req.MuxNocomp),
			)
		}

		if req.EnableKcp {
			opts = append(opts,
				well_known_option.WithKcpPreset(req.KcpPreset),
				well_known_option.WithKcpCrypt(req.KcpCrypt),
				well_known_option.WithKcpKey(req.KcpKey),
				well_known_option.WithKcpMtu(req.KcpMtu),
				well_known_option.WithKcpSndwnd(req.KcpSndwnd),
				well_known_option.WithKcpRecvwnd(req.KcpRcvwnd),
				well_known_option.WithKcpDataShard(req.KcpDataShard),
				well_known_option.WithKcpParityShard(req.KcpParityShard),
			)
		}
	}

	tp, err := mp.Teleport(ctx, target, req.SourceNetwork, req.SourceAddress, req.SinkNetwork, req.SinkAddress, opts...)
	if err != nil {
		logger.WithError(err).Errorf("failed to teleport")
		return
	}

	logger.Infof("teleport")

	return ViewTeleportation(tp), nil
}
