package meepo_core

import (
	"context"
	"encoding/hex"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
	"github.com/PeerXu/meepo/pkg/meepo/transport"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

func (mp *Meepo) newOnNewTransportRequest() any { return &crypto_core.Packet{} }

func (mp *Meepo) hdrOnNewTransport(ctx context.Context, req any) (any, error) {
	return mp.onNewTransport(req.(*crypto_core.Packet))
}

func (mp *Meepo) onNewTransport(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
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
		answer, err = mp.forwardNewTransportRequestToClosestTrackers(dstAddr, in)
		if err != nil {
			logger.WithError(err).Debugf("failed to forward new transport request to closest trackers")
			return
		}
		logger.Tracef("forward new transport request")
		return
	}

	buf, err := mp.cryptor.Decrypt(in)
	if err != nil {
		logger.WithError(err).Debugf("failed to decrypt packet")
		return
	}

	var req tracker_interface.NewTransportRequest
	if err = mp.unmarshaler.Unmarshal(buf, &req); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal buffer")
		return
	}

	pc, err := mp.newPeerConnection()
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return
	}

	done := make(chan struct{})
	defer close(done)
	var er error

	mp.transportsMtx.Lock()

	if _, found := mp.transports[srcAddr]; found {
		defer mp.transportsMtx.Unlock()
		pc.Close()
		err = ErrTransportFoundFn(srcAddr.String())
		logger.WithError(err).Debugf("transport existed")
		return
	}

	opts := []NewTransportOption{
		well_known_option.WithLogger(mp.GetLogger()),
		well_known_option.WithAddr(srcAddr),
		transport_webrtc.WithOffer(req.Offer),
		transport_webrtc.WithSession(req.Session),
		well_known_option.WithPeerConnection(pc),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		marshaler.WithMarshaler(mp.marshaler),
		marshaler.WithUnmarshaler(mp.unmarshaler),
		transport_webrtc.WithGatherDoneFunc(func(_ transport_webrtc.Session, _answer webrtc.SessionDescription, _err error) {
			answer = _answer
			er = _err
			done <- struct{}{}
		}),
		transport_core.WithBeforeNewChannelHook(mp.beforeNewChannelHook),
		transport_core.WithOnTransportCloseFunc(mp.onRemoveWebrtcTransport),
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

	t, err := transport.NewTransport(transport_webrtc.TRANSPORT_WEBRTC_SINK, opts...)
	if err != nil {
		defer mp.transportsMtx.Unlock()
		logger.WithError(err).Debugf("failed to new transport")
		return
	}

	mp.onAddWebrtcTransportNL(t) // nolint:errcheck

	mp.transportsMtx.Unlock()

	<-done
	if er != nil {
		mp.transportsMtx.Lock()
		defer mp.transportsMtx.Unlock()
		defer t.Close(mp.context())
		delete(mp.transports, srcAddr)

		err = er
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	logger.Tracef("on new transport")
	return
}

func (mp *Meepo) forwardNewTransportRequestToClosestTrackers(dstAddr addr.Addr, in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":     "forwardNewTransportRequestToClosestTrackers",
		"destination": dstAddr.String(),
		"nonce":       hex.EncodeToString(in.Nonce),
	})

	tks, err := mp.getClosestTrackers(dstAddr)
	if err != nil {
		logger.WithError(err).Debugf("failed to get closest trackers")
		return
	}
	logger = logger.WithField("trackers.length", len(tks))

	if len(tks) == 0 {
		err = ErrNoAvailableTrackers
		logger.WithError(err).Debugf("no available trackers")
		return
	}

	done := make(chan struct{})
	answers := make(chan webrtc.SessionDescription)
	errs := make(chan error)
	defer func() {
		close(done)
		close(answers)
		close(errs)
	}()

	for _, tk := range tks {
		go func(tk Tracker) {
			logger := logger.WithField("tracker", tk.Addr().String())
			_answer, _err := tk.NewTransport(in)
			select {
			case <-done:
				logger.Tracef("forward already done")
				return
			default:
			}

			if _err != nil {
				logger.WithError(_err).Tracef("failed to new transport by tracker")
				errs <- _err
				return
			}
			answers <- _answer
			logger.Tracef("new transport by tracker")
		}(tk)
	}

	for i := 0; i < len(tks); i++ {
		select {
		case answer = <-answers:
			logger.Tracef("forward done")
			return answer, nil
		case err = <-errs:
		}
	}

	logger.WithError(err).Debugf("failed to forward")
	return
}
