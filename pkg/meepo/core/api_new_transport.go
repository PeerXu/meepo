package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPINewTransport(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.NewTransportRequest)

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrAPINewTransport",
		"target":  req.Target,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return nil, err
	}

	var opts []NewTransportOption

	if req.Manual {
		opts = []sdk_interface.NewTransportOption{
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

	t, err := mp.NewTransport(ctx, target, opts...)
	if err != nil {
		logger.WithError(err).Errorf("failed to new transport")
		return nil, err
	}

	logger.Infof("new transport")

	return ViewTransport(t), nil
}
