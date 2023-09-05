package meepo_core

import (
	"context"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
	"github.com/PeerXu/meepo/pkg/meepo/transport"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

func (mp *Meepo) hdrOnNewTransport(ctx context.Context, req any) (any, error) {
	return mp.onNewTransport(req.(*crypto_core.Packet))
}

func (mp *Meepo) onNewTransport(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	var t Transport

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onNewTransport",
	})

	srcAddr, err := addr.FromBytesWithoutMagicCode(in.Source)
	if err != nil {
		logger.WithError(err).Debugf("invalid source address")
		return
	}
	logger = logger.WithField("source", srcAddr.String())

	dstAddr, err := addr.FromBytesWithoutMagicCode(in.Destination)
	if err != nil {
		logger.WithError(err).Debugf("invalid destination address")
		return
	}
	logger = logger.WithField("destination", dstAddr.String())

	if err = mp.signer.Verify(in); err != nil {
		logger.WithError(err).Debugf("failed to verify packet")
		return
	}

	if !mp.Addr().Equal(dstAddr) {
		answer, err = mp.forwardNewTransportRequest(dstAddr, in)
		if err != nil {
			logger.WithError(err).Debugf("failed to forward new transport request to closest trackers")
			return
		}
		logger.Tracef("forward new transport request")
		return
	}

	var req tracker_interface.NewTransportRequest
	if err = mp.decryptMessage(in, &req); err != nil {
		logger.WithError(err).Debugf("failed to decrypt message")
		return
	}

	done := make(chan struct{})
	defer close(done)
	var er error

	mp.transportsMtx.Lock()

	if _, found := mp.transports[srcAddr]; found {
		defer mp.transportsMtx.Unlock()
		err = ErrTransportFoundFn(srcAddr.String())
		logger.WithError(err).Debugf("transport existed")
		return
	}

	opts := []NewTransportOption{
		well_known_option.WithLogger(mp.GetLogger()),
		well_known_option.WithAddr(srcAddr),
		transport_webrtc.WithOffer(req.Offer),
		transport_webrtc.WithSession(req.Session),
		transport_webrtc.WithNewPeerConnectionFunc(mp.newPeerConnection),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		marshaler.WithMarshaler(mp.marshaler),
		marshaler.WithUnmarshaler(mp.unmarshaler),
		transport_webrtc.WithGatherDoneOnNewFunc(func(_ transport_webrtc.Session, _answer webrtc.SessionDescription, _err error) {
			answer = _answer
			er = _err
			done <- struct{}{}
		}),
		transport_webrtc.WithGatherFunc(mp.genGatherFunc(srcAddr)),
		transport_core.WithAfterNewTransportHook(func(t meepo_interface.Transport, opts ...transport_core.HookOption) {
			mp.onAddWebrtcTransportNL(t)
			mp.emitTransportActionNew(t)
		}),
		transport_core.WithAfterCloseTransportHook(func(t meepo_interface.Transport, opts ...transport_core.HookOption) {
			mp.onRemoveWebrtcTransport(t)
			mp.emitTransportActionClose(t)
		}),
		transport_core.WithBeforeNewChannelHook(func(network, address string, opts ...transport_core.HookOption) error {
			return mp.beforeNewChannelHook(t, network, address, opts...)
		}),
		transport_core.WithAfterNewChannelHook(func(c meepo_interface.Channel, opts ...transport_core.HookOption) {
			mp.emitChannelActionNew(c)
		}),
		transport_core.WithAfterCloseChannelHook(func(c meepo_interface.Channel, opts ...transport_core.HookOption) {
			mp.emitChannelActionClose(c)
		}),
		transport_core.WithOnTransportStateChangeFunc(func(t meepo_interface.Transport) {
			mp.emitTransportStateChange(t)
		}),
		transport_core.WithOnChannelStateChangeFunc(func(c meepo_interface.Channel) {
			mp.emitChannelStateChange(srcAddr, c)
		}),
		transport_core.WithOnTransportReadyFunc(mp.onReadyWebrtcTransport),
		well_known_option.WithEnableMux(req.EnableMux),
		well_known_option.WithEnableKcp(req.EnableKcp),
	}
	if req.EnableMux {
		opts = append(opts,
			transport_webrtc.WithMuxLabel(req.MuxLabel),
			well_known_option.WithMuxVer(req.MuxVer),
			well_known_option.WithMuxBuf(req.MuxBuf),
			well_known_option.WithMuxStreamBuf(req.MuxStreamBuf),
			well_known_option.WithMuxNocomp(req.MuxNocomp),
		)
	}

	if req.EnableKcp {
		opts = append(opts,
			transport_webrtc.WithKcpLabel(req.KcpLabel),
			well_known_option.WithKcpPreset(req.KcpPreset),
			well_known_option.WithKcpCrypt(req.KcpCrypt),
			well_known_option.WithKcpKey(req.KcpKey),
			well_known_option.WithKcpMtu(req.KcpMtu),
			well_known_option.WithKcpSndwnd(req.KcpSndwnd),
			well_known_option.WithKcpRecvwnd(req.KcpRcvwnd),
			well_known_option.WithKcpDataShard(req.DataShard),
			well_known_option.WithKcpParityShard(req.ParityShard),
		)
	}

	t, err = transport.NewTransport(transport_webrtc.TRANSPORT_WEBRTC_SINK, opts...)
	if err != nil {
		defer mp.transportsMtx.Unlock()
		logger.WithError(err).Debugf("failed to new transport")
		return
	}

	mp.transportsMtx.Unlock()

	<-done
	if er != nil {
		mp.transportsMtx.Lock()
		defer mp.transportsMtx.Unlock()
		defer t.Close(mp.context())
		mp.removeTransportNL(srcAddr)
		err = er
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	logger.Tracef("on new transport")
	return
}

func (mp *Meepo) forwardNewTransportRequest(dstAddr addr.Addr, in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "forwardNewTransportRequest",
		"target":  dstAddr.String(),
	})

	out, err := mp.forwardRequest(mp.context(), dstAddr, in, func(tk Tracker, in *crypto_core.Packet) (any, error) { return tk.NewTransport(in) }, mp.getClosestTrackers, logger)
	if err != nil {
		return
	}

	answer = out.(webrtc.SessionDescription)

	return
}
