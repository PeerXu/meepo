package meepo_core

import (
	"context"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/dialer"
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
	"github.com/PeerXu/meepo/pkg/meepo/transport"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

func (mp *Meepo) NewTransport(ctx context.Context, target Addr, opts ...NewTransportOption) (Transport, error) {
	var name string
	var err error

	o := option.ApplyWithDefault(mp.defaultNewTransportOptions(), opts...)
	gtkFn, _ := GetGetTrackersFunc(o)
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "NewTransport",
		"target":  target,
	})

	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()

	if _, found := mp.transports[target]; found {
		err = ErrTransportExistFn(target.String())
		logger.WithError(err).Debugf("tansport exist")
		return nil, err
	}

	var onAddTransport func(Transport) error
	var onReadyTransport func(Transport) error
	var onRemoveTransport func(Transport) error
	switch target {
	case mp.Addr():
		name = "pipe"

		onAddTransport = mp.onAddPipeTransportNL
		onRemoveTransport = mp.onRemovePipeTransport
		onReadyTransport = func(Transport) error { return nil }
	default:
		name = "webrtc/source"
		var gatherOpt gatherOption

		gatherOpt.EnableMux, _ = well_known_option.GetEnableMux(o)
		if gatherOpt.EnableMux {
			gatherOpt.MuxLabel = mp.newLabel("mux")
			gatherOpt.MuxVer, _ = well_known_option.GetMuxVer(o)
			gatherOpt.MuxBuf, _ = well_known_option.GetMuxBuf(o)
			gatherOpt.MuxStreamBuf, _ = well_known_option.GetMuxStreamBuf(o)
			gatherOpt.MuxNocomp, _ = well_known_option.GetMuxNocomp(o)

			opts = append(opts,
				transport_webrtc.WithMuxLabel(gatherOpt.MuxLabel),
				well_known_option.WithMuxVer(gatherOpt.MuxVer),
				well_known_option.WithMuxBuf(gatherOpt.MuxBuf),
				well_known_option.WithMuxStreamBuf(gatherOpt.MuxStreamBuf),
				well_known_option.WithMuxNocomp(gatherOpt.MuxNocomp),
			)
		}

		gatherOpt.EnableKcp, _ = well_known_option.GetEnableKcp(o)
		if gatherOpt.EnableKcp {
			gatherOpt.KcpLabel = mp.newLabel("kcp")
			gatherOpt.KcpPreset, _ = well_known_option.GetKcpPreset(o)
			gatherOpt.KcpCrypt, _ = well_known_option.GetKcpCrypt(o)
			gatherOpt.KcpKey, _ = well_known_option.GetKcpKey(o)
			gatherOpt.KcpMtu, _ = well_known_option.GetKcpMtu(o)
			gatherOpt.KcpSndwnd, _ = well_known_option.GetKcpSndwnd(o)
			gatherOpt.KcpRcvwnd, _ = well_known_option.GetKcpRcvwnd(o)
			gatherOpt.KcpDataShard, _ = well_known_option.GetKcpDataShard(o)
			gatherOpt.KcpParityShard, _ = well_known_option.GetKcpParityShard(o)

			opts = append(opts,
				transport_webrtc.WithKcpLabel(gatherOpt.KcpLabel),
				well_known_option.WithKcpPreset(gatherOpt.KcpPreset),
				well_known_option.WithKcpCrypt(gatherOpt.KcpCrypt),
				well_known_option.WithKcpKey(gatherOpt.KcpKey),
				well_known_option.WithKcpMtu(gatherOpt.KcpMtu),
				well_known_option.WithKcpSndwnd(gatherOpt.KcpSndwnd),
				well_known_option.WithKcpRecvwnd(gatherOpt.KcpRcvwnd),
				well_known_option.WithKcpDataShard(gatherOpt.KcpDataShard),
				well_known_option.WithKcpParityShard(gatherOpt.KcpParityShard),
			)
		}
		opts = append(opts,
			well_known_option.WithEnableMux(gatherOpt.EnableMux),
			well_known_option.WithEnableKcp(gatherOpt.EnableKcp),
		)

		pc, err := mp.newPeerConnection()
		if err != nil {
			logger.WithError(err).Debugf("failed to new peer connection")
			return nil, err
		}

		opts = append(opts,
			transport_webrtc.WithGatherFunc(mp.gatherFunc(target, gtkFn, gatherOpt)),
			well_known_option.WithPeerConnection(pc),
			transport_core.WithBeforeNewChannelHook(func(t Transport, network, address string) error {
				return mp.permit(t.Addr().String(), network, address)
			}),
		)

		onAddTransport = mp.onAddWebrtcTransportNL
		onRemoveTransport = mp.onRemoveWebrtcTransport
		onReadyTransport = mp.onReadyWebrtcTransport
	}
	logger = logger.WithField("name", name)
	opts = append(opts,
		dialer.WithDialer(dialer.GetGlobalDialer()),
		well_known_option.WithAddr(target),
		well_known_option.WithLogger(mp.GetRawLogger()),
		transport_core.WithOnTransportCloseFunc(onRemoveTransport),
		transport_core.WithOnTransportReadyFunc(onReadyTransport),
		marshaler.WithMarshaler(mp.marshaler),
		marshaler.WithUnmarshaler(mp.unmarshaler),
	)

	t, err := transport.NewTransport(name, opts...)
	if err != nil {
		logger.WithError(err).Debugf("failed to new transport")
		return nil, err
	}

	onAddTransport(t) // nolint:errcheck

	logger.Tracef("new transport")

	return t, nil
}

func (mp *Meepo) newNewTransportRequest(target addr.Addr, offer webrtc.SessionDescription, opt gatherOption) (*crypto_core.Packet, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "newNewTransportRequest",
		"target":  target.String(),
	})

	req := &tracker_interface.NewTransportRequest{
		Offer: offer,

		EnableMux:    opt.EnableMux,
		MuxLabel:     opt.MuxLabel,
		MuxVer:       opt.MuxVer,
		MuxBuf:       opt.MuxBuf,
		MuxStreamBuf: opt.MuxStreamBuf,
		MuxNocomp:    opt.MuxNocomp,

		EnableKcp:   opt.EnableKcp,
		KcpLabel:    opt.KcpLabel,
		KcpPreset:   opt.KcpPreset,
		KcpCrypt:    opt.KcpCrypt,
		KcpKey:      opt.KcpKey,
		KcpMtu:      opt.KcpMtu,
		KcpSndwnd:   opt.KcpSndwnd,
		KcpRcvwnd:   opt.KcpRcvwnd,
		DataShard:   opt.KcpDataShard,
		ParityShard: opt.KcpParityShard,
	}

	buf, err := mp.marshaler.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal offer to plaintext")
		return nil, err
	}

	out, err := mp.cryptor.Encrypt(target.Bytes(), buf)
	if err != nil {
		logger.WithError(err).Debugf("failed to encrypt plaintext to packet")
		return nil, err
	}

	if err = mp.signer.Sign(out); err != nil {
		logger.WithError(err).Debugf("failed to sign packet")
		return nil, err
	}

	logger.Tracef("new NewTransport request")

	return out, nil
}
