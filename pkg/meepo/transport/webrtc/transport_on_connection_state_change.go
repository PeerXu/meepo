package transport_webrtc

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) onSinkConnectionStateChange(sess Session) func(webrtc.PeerConnectionState) {
	return func(s webrtc.PeerConnectionState) {
		logger := t.GetLogger().WithFields(logging.Fields{
			"#method": "onSinkConnectionStateChange",
			"session": sess.String(),
			"state":   s.String(),
		})

		switch s {
		case webrtc.PeerConnectionStateNew:

		case webrtc.PeerConnectionStateConnecting:

		case webrtc.PeerConnectionStateConnected:
			atomic.StoreInt64(&t.stat.failedSinkConnections, 0)
			if t.ensureUniqueConnectedPeerConnection(sess) {
				logger.Tracef("peer connection is unique")

				go t.tryNewSysDataChannel(sess)
				if t.enableMux {
					go t.tryNewMuxDataChannel()
				}
				if t.enableKcp {
					go t.tryNewKcpDataChannel()
				}
			}
		case webrtc.PeerConnectionStateDisconnected:

		case webrtc.PeerConnectionStateFailed:
			atomic.AddInt64(&t.stat.failedSinkConnections, 1)

			pc, err := t.loadPeerConnection(sess)
			if err != nil {
				logger.WithError(err).Debugf("failed to get load peer connection")
			} else {
				go logger.WithError(pc.Close()).Tracef("close peer connection")
			}
		case webrtc.PeerConnectionStateClosed:
			t.unregisterPeerConnection(sess)
			if t.isClosable() {
				go t.Close(t.context())
			}
		}

		logger.Tracef("peer connection state changed")
	}
}

func (t *WebrtcTransport) onSourceConnectionStateChange(sess Session) func(webrtc.PeerConnectionState) {
	return func(s webrtc.PeerConnectionState) {
		logger := t.GetLogger().WithFields(logging.Fields{
			"#method": "onSourceConnectionStateChange",
			"session": sess.String(),
			"state":   s.String(),
		})

		switch s {
		case webrtc.PeerConnectionStateNew:

		case webrtc.PeerConnectionStateConnecting:

		case webrtc.PeerConnectionStateConnected:
			atomic.StoreInt64(&t.stat.failedSourceConnections, 0)
			if t.ensureUniqueConnectedPeerConnection(sess) {
				logger.Tracef("peer connection is unique")
			}
		case webrtc.PeerConnectionStateDisconnected:

		case webrtc.PeerConnectionStateFailed:
			atomic.AddInt64(&t.stat.failedSourceConnections, 1)

			pc, err := t.loadPeerConnection(sess)
			if err != nil {
				logger.WithError(err).Debugf("failed to get load peer connection")
			} else {
				go logger.WithError(pc.Close()).Tracef("close peer connection")
			}
		case webrtc.PeerConnectionStateClosed:
			t.unregisterPeerConnection(sess)
			if t.isClosable() {
				go t.Close(t.context())
			}
		}

		logger.Tracef("peer connection state changed")
	}
}

func (t *WebrtcTransport) tryNewSysDataChannel(sess Session) {
	var err error

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "tryNewSysDataChannel",
		"session": sess.String(),
	})

	t.sysMtx.Lock()
	defer t.sysMtx.Unlock()

	pc, err := t.loadPeerConnection(sess)
	if err != nil {
		logger.WithError(err).Debugf("failed to load peer connection")
		return
	}

	sdc, err := t.loadSystemDataChannel(sess)
	if err == nil {
		if t.isSysDataChannelAlive(sess) {
			logger.WithField("label", sdc.Label()).Tracef("sys data channel still alive, skip")
			return
		}

		logger.WithFields(logging.Fields{
			"old": sdc.Label(),
		}).WithError(sdc.Close()).Debugf("close old system data channel")
		t.unregisterSystemDataChannel(sess)
		t.unregisterSystemReadWriteCloser(sess)
	} else {
		if !errors.Is(err, ErrDataChannelNotFound) {
			logger.WithError(err).Debugf("failed to get system data channel")
			return
		}
	}

	ordered := !t.enableKcp
	sdc, err = pc.CreateDataChannel(t.nextDataChannelLabel("sys"), &webrtc.DataChannelInit{Ordered: &ordered})
	if err != nil {
		defer t.Close(t.context())
		logger.WithError(err).Debugf("failed to new sys data channel")
		return
	}
	logger = logger.WithField("label", sdc.Label())
	t.registerSystemDataChannel(sess, sdc)

	sdc.OnOpen(func() { t.onSysDataChannelOpen(sess, sdc) })

	logger.Tracef("new sys data channel")
}

func (t *WebrtcTransport) isSysDataChannelAlive(sess Session) bool {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "isSysDataChannelAlive",
		"session": sess,
	})

	ctx := context.WithValue(t.context(), OPTION_SESSION, sess)
	ttl, err := t.ping(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to ping")
		return false
	}

	logger.WithField("ttl", ttl).Tracef("sys data channel alive")

	return true
}

func (t *WebrtcTransport) nextDataChannelLabel(scope string) string {
	return fmt.Sprintf("%s#%016x", scope, t.randSrc.Int63())
}

func (t *WebrtcTransport) tryNewMuxDataChannel() {
	t.tryNewDataChannelTpl("mux", t.muxLabel, t.muxMtx, t.GetLogger().WithFields(logging.Fields{
		"#method": "tryNewMuxDataChannel",
		"label":   t.muxLabel,
	}), &t.muxDataChannel, t.onMuxDataChannelOpen, true)()
}

func (t *WebrtcTransport) tryNewKcpDataChannel() {
	t.tryNewDataChannelTpl("kcp", t.kcpLabel, t.kcpMtx, t.GetLogger().WithFields(logging.Fields{
		"#method": "tryNewKcpDataChannel",
		"label":   t.kcpLabel,
	}), &t.kcpDataChannel, t.onKcpDataChannelOpen, false)()
}

func (t *WebrtcTransport) tryNewDataChannelTpl(name string, label string, mtx lock.Locker, logger logging.Logger, dcPtr **webrtc.DataChannel, onOpen func(dc *webrtc.DataChannel), ordered bool) func() {
	return func() {
		var err error

		dc := *dcPtr

		mtx.Lock()
		defer mtx.Unlock()

		if isDataChannelAlive(dc) {
			logger.Tracef("%s data channel still alive, skip", name)
			return
		}

		if dc != nil {
			logger.WithError(dc.Close()).Debugf("close old %s data channel", name)
			*dcPtr = nil
		}

		pc, err := t.loadRandomPeerConnection()
		if err != nil {
			// TODO: dont close transport when failed to load pc
			defer t.Close(t.context())
			logger.WithError(err).Debugf("failed to load peer connection")
			return
		}

		ndc, err := pc.CreateDataChannel(label, &webrtc.DataChannelInit{Ordered: &ordered})
		if err != nil {
			defer t.Close(t.context())
			logger.WithError(err).Debugf("failed to new %s data channel", name)
			return
		}

		ndc.OnOpen(func() { onOpen(ndc) })
		logger.Tracef("new %s data channel", name)

		*dcPtr = ndc
	}
}

func isDataChannelAlive(dc *webrtc.DataChannel) bool {
	if dc == nil {
		return false
	}
	return dc.ReadyState() == webrtc.DataChannelStateOpen
}
