package meepo

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
)

func newTeleportOption() objx.Map {
	return objx.New(map[string]interface{}{})
}

func (mp *Meepo) Teleport(id string, remote net.Addr, opts ...TeleportOption) (net.Addr, error) {
	var ok bool
	var err error

	logger := mp.getLogger().WithField("#method", "Teleport")

	o := newTeleportOption()
	for _, opt := range opts {
		opt(o)
	}

	var local net.Addr
	if local, ok = o.Get("local").Inter().(net.Addr); ok {
		if local, ok = isListenableAddr(local); !ok {
			err = NotListenableAddressError(local)
			logger.WithError(err).Debugf("not listenable address")
			return nil, err
		}
	} else {
		local = getListenableAddr()
	}

	logger = logger.WithFields(logrus.Fields{
		"network": local.Network(),
		"address": local.String(),
	})

	lis, err := net.Listen(local.Network(), local.String())
	if err != nil {
		logger.WithError(err).Debugf("failed to listen")
		return nil, err
	}
	logger.Tracef("listen")

	go mp.doTeleport(id, remote, lis)

	return local, nil
}

func (mp *Meepo) doTeleport(id string, remote net.Addr, lis net.Listener) {
	logger := mp.getLogger().WithField("#method", "doTeleport")

	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields{
			"id":    id,
			"laddr": lis.Addr(),
			"raddr": remote,
		})
	}

	defer func() { logger.WithError(lis.Close()).Tracef("listen closed") }()

	conn, err := lis.Accept()
	if err != nil {
		logger.WithError(err).Debugf("failed to accept")
		return
	}
	defer func() { logger.WithError(conn.Close()).Tracef("connection closed") }()
	logger.Tracef("accept")

	gatherer, err := mp.rtc.NewICEGatherer(mp.iceGatherOptions())
	if err != nil {
		logger.WithError(err).Debugf("failed to new ice gatherer")
		return
	}
	logger.Tracef("new ice gathererer")

	ice := mp.rtc.NewICETransport(gatherer)
	logger.Tracef("new ice transport")

	dtls, err := mp.rtc.NewDTLSTransport(ice, nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to new ice transport")
		return
	}
	logger.Tracef("new dtls transport")

	sctp := mp.rtc.NewSCTPTransport(dtls)
	logger.Tracef("new sctp transport")

	gatherFinished := make(chan struct{})
	gatherer.OnLocalCandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			close(gatherFinished)
		}
	})
	if err = gatherer.Gather(); err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return
	}
	logger.Tracef("gather started")

	select {
	case <-gatherFinished:
		logger.Tracef("gather done")
	case <-time.After(cast.ToDuration(mp.opt.Get("gatherTimeout").Inter())):
		logger.Debugf("gather timeout")
		return
	}

	iceCandidates, err := gatherer.GetLocalCandidates()
	if err != nil {
		logger.WithError(err).Debugf("failed to get local candidates")
		return
	}
	logger.Tracef("get local gatherer candidates")

	iceParams, err := gatherer.GetLocalParameters()
	if err != nil {
		logger.WithError(err).Debugf("failed to get gatherer local parameters")
		return
	}
	logger.Tracef("get local gatherer parameters")

	dtlsParams, err := dtls.GetLocalParameters()
	if err != nil {
		logger.WithError(err).Debugf("failed to get dtls local parameters")
		return
	}
	logger.Tracef("get local dtls parameters")

	sctpCapbilities := sctp.GetCapabilities()

	src := &signaling.Descriptor{
		ID: mp.ID(),
		Signal: &signaling.Signal{
			ICECandidates:    iceCandidates,
			ICEParameters:    iceParams,
			DTLSParameters:   dtlsParams,
			SCTPCapabilities: sctpCapbilities,
		},
		UserData: map[string]interface{}{
			"network": remote.Network(),
			"address": remote.String(),
		},
	}

	dst, err := mp.se.Wire(&signaling.Descriptor{ID: id}, src)
	if err != nil {
		logger.WithError(err).Debugf("failed to wire")
		return
	}

	var eg ImmediatelyErrorGroup

	if err = ice.SetRemoteCandidates(dst.Signal.ICECandidates); err != nil {
		logger.WithError(err).Debugf("failed to set remote ice candidates")
		return
	}
	logger.Tracef("set remote ice candidates")

	iceRole := webrtc.ICERoleControlling
	if err = ice.Start(nil, dst.Signal.ICEParameters, &iceRole); err != nil {
		logger.WithError(err).Debugf("failed to start ice transport")
		return
	}
	ice.OnConnectionStateChange(func(s webrtc.ICETransportState) {
		logger.WithField("state", s).Debugf("ice transport state changed")
		if s == webrtc.ICETransportStateFailed {
			eg.Go(func() error { return fmt.Errorf("ICETransportStateFailed") })
		}
	})
	defer func() { logger.WithError(ice.Stop()).Tracef("ice closed") }()
	logger.Tracef("start ice transport")

	if err = dtls.Start(dst.Signal.DTLSParameters); err != nil {
		logger.WithError(err).Debugf("failed to start dtls transport")
		return
	}
	dtls.OnStateChange(func(s webrtc.DTLSTransportState) {
		logger.WithField("state", s).Debugf("dtls transport state changed")
		if s == webrtc.DTLSTransportStateFailed {
			eg.Go(func() error { return fmt.Errorf("DTLSTransportStateFailed") })
		}
	})
	defer func() { logger.WithError(dtls.Stop()).Tracef("dtls closed") }()
	logger.Tracef("start dtls transport")

	if err = sctp.Start(dst.Signal.SCTPCapabilities); err != nil {
		logger.WithError(err).Debugf("failed to start sctp transport")
		return
	}
	sctp.OnError(func(err error) {
		logger.WithError(err).Debugf("sctp error")
		eg.Go(func() error { return err })
	})
	defer func() { logger.WithError(sctp.Stop()).Tracef("sctp closed") }()
	logger.Tracef("start sctp transport")

	var dcID uint16 = 1
	dcParams := &webrtc.DataChannelParameters{
		Label: "meepo.teleport",
		ID:    &dcID,
	}

	channel, err := mp.rtc.NewDataChannel(sctp, dcParams)
	if err != nil {
		logger.WithError(err).Debugf("failed to new data channel")
		return
	}
	channel.OnError(func(err error) {
		logger.WithError(err).Debugf("channel error")
		eg.Go(func() error { return err })
	})
	defer func() { logger.WithError(channel.Close()).Tracef("channel closed") }()
	logger.Tracef("new data channel")

	var wg sync.WaitGroup
	wg.Add(1)
	channel.OnOpen(func() {
		chConn, err := channel.Detach()
		if err != nil {
			logger.WithError(err).Debugf("failed to detach data channel")
			return
		}

		eg.Go(func() error {
			_, err := io.Copy(conn, chConn)
			logger.WithError(err).Tracef("conn<-chConn closed")
			return err
		})
		eg.Go(func() error {
			_, err := io.Copy(chConn, conn)
			logger.WithError(err).Tracef("conn->chConn closed")
			return err
		})

		wg.Done()
		logger.Tracef("data channel opened")
	})
	wg.Wait()

	logger.WithError(eg.Wait()).Debugf("done")
}

func (mp *Meepo) onTeleport(src *signaling.Descriptor) (*signaling.Descriptor, error) {
	logger := mp.getLogger().WithField("#method", "onTeleport")

	var eg ImmediatelyErrorGroup
	var outerDialCloser func() error

	ud := objx.New(src.UserData)
	lnet, laddr := ud.Get("network").String(), ud.Get("address").String()
	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields{
			"network": lnet,
			"address": laddr,
		})
	}

	addr, err := net.ResolveTCPAddr(lnet, laddr)
	if err != nil {
		logger.WithError(err).Debugf("failed to resolve tcp addr")
		return nil, err
	}

	dial, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		logger.WithError(err).Debugf("failed to dial")
		return nil, err
	}
	outerDialCloser = dial.Close
	defer func() { outerDialCloser() }()
	logger.Tracef("dial")

	gatherer, err := mp.rtc.NewICEGatherer(mp.iceGatherOptions())
	if err != nil {
		logger.WithError(err).Debugf("failed to new ice gatherer")
		return nil, err
	}
	logger.Tracef("new ice gatherer")

	ice := mp.rtc.NewICETransport(gatherer)
	logger.Tracef("new ice transport")

	dtls, err := mp.rtc.NewDTLSTransport(ice, nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to new dtls transport")
		return nil, err
	}
	logger.Tracef("new dtls transport")

	sctp := mp.rtc.NewSCTPTransport(dtls)
	logger.Tracef("new sctp transport")

	var wg sync.WaitGroup
	wg.Add(1)
	sctp.OnDataChannel(func(channel *webrtc.DataChannel) {
		channel.OnOpen(func() {
			chConn, err := channel.Detach()
			if err != nil {
				logger.WithError(err).Debugf("failed to detach data channel")
				return
			}

			eg.Go(func() error {
				_, err = io.Copy(dial, chConn)
				logger.WithError(err).Debugf("dial<-chConn closed")
				return err
			})
			eg.Go(func() error {
				_, err = io.Copy(chConn, dial)
				logger.WithError(err).Debugf("dial->chConn closed")
				return err
			})

			wg.Done()
			logger.Tracef("data channel opened")
		})
		channel.OnError(func(err error) {
			logger.WithError(err).Debugf("channel error")
			eg.Go(func() error { return err })
		})
	})

	gatherFinished := make(chan struct{})
	gatherer.OnLocalCandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			close(gatherFinished)
		}
	})
	if err = gatherer.Gather(); err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return nil, err
	}
	logger.Tracef("gather started")

	select {
	case <-gatherFinished:
		logger.Tracef("gather done")
	case <-time.After(cast.ToDuration(mp.opt.Get("gatherTimeout").Inter())):
		logger.Debugf("gather timeout")
		return nil, GatherTimeoutError
	}

	iceCandidates, err := gatherer.GetLocalCandidates()
	if err != nil {
		logger.WithError(err).Debugf("failed to get gatherer local candidates")
		return nil, err
	}
	logger.Tracef("get local gatherer candidates")

	iceParams, err := gatherer.GetLocalParameters()
	if err != nil {
		logger.WithError(err).Debugf("failed to get gatherer local parameters")
		return nil, err
	}
	logger.Tracef("get local gatherer parameters")

	dtlsParams, err := dtls.GetLocalParameters()
	if err != nil {
		logger.WithError(err).Debugf("failed to get dtls local parameters")
		return nil, err
	}
	logger.Tracef("get local dtls parameters")

	sctpCapbilities := sctp.GetCapabilities()

	dst := &signaling.Descriptor{
		ID: mp.ID(),
		Signal: &signaling.Signal{
			ICECandidates:    iceCandidates,
			ICEParameters:    iceParams,
			DTLSParameters:   dtlsParams,
			SCTPCapabilities: sctpCapbilities,
		},
	}

	go func() {
		defer func() { logger.WithError(dial.Close()).Tracef("connection closed") }()

		if err = ice.SetRemoteCandidates(src.Signal.ICECandidates); err != nil {
			logger.WithError(err).Debugf("failed to set remote candidates")
			return
		}
		logger.Tracef("set remote ice candidates")

		iceRole := webrtc.ICERoleControlled
		if err = ice.Start(nil, src.Signal.ICEParameters, &iceRole); err != nil {
			logger.WithError(err).Debugf("failed to start ice transport")
			return
		}
		ice.OnConnectionStateChange(func(s webrtc.ICETransportState) {
			logger.WithField("state", s).Tracef("ice transport state changed")
			if s == webrtc.ICETransportStateFailed {
				eg.Go(func() error { return fmt.Errorf("ICETransportStateFailed") })
			}
		})
		defer func() { logger.WithError(ice.Stop()).Tracef("ice closed") }()
		logger.Tracef("start ice transport")

		if err = dtls.Start(src.Signal.DTLSParameters); err != nil {
			logger.WithError(err).Debugf("failed to start dtls transport")
			return
		}
		dtls.OnStateChange(func(s webrtc.DTLSTransportState) {
			logger.WithField("state", s).Tracef("dtls transport state changed")
			if s == webrtc.DTLSTransportStateFailed {
				eg.Go(func() error { return fmt.Errorf("DTLSTransportStateFailed") })
			}
		})
		defer func() { logger.WithError(dtls.Stop()).Tracef("dtls closed") }()
		logger.Tracef("start dtls transport")

		if err = sctp.Start(src.Signal.SCTPCapabilities); err != nil {
			logger.WithError(err).Debugf("failed to start sctp transport")
			return
		}
		sctp.OnError(func(err error) {
			logger.WithError(err).Debugf("sctp error")
			eg.Go(func() error { return err })
		})
		defer func() { logger.WithError(sctp.Stop()).Tracef("sctp closed") }()
		logger.Tracef("start sctp transport")

		wg.Wait()

		logger.WithError(eg.Wait()).Debugf("done")
	}()

	outerDialCloser = func() error { return nil }

	return dst, nil
}
