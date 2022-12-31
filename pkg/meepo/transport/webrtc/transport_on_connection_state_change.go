package transport_webrtc

import (
	"fmt"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/lock"
)

func (t *WebrtcTransport) onSinkConnectionStateChange(s webrtc.PeerConnectionState) {
	logger := t.GetLogger().WithField("#method", "onSinkConnectionStateChange")

	switch s {
	case webrtc.PeerConnectionStateNew,
		webrtc.PeerConnectionStateConnecting:
	case webrtc.PeerConnectionStateConnected:
		go t.tryNewSysDataChannel()
		if t.enableMux {
			go t.tryNewMuxDataChannel()
		}
		if t.enableKcp {
			go t.tryNewKcpDataChannel()
		}
	case webrtc.PeerConnectionStateDisconnected:
	case webrtc.PeerConnectionStateFailed,
		webrtc.PeerConnectionStateClosed:
		go t.Close(t.context())
	}

	logger.WithField("state", s.String()).Tracef("peer connection state changed")
}

func (t *WebrtcTransport) onSourceConnectionStateChange(s webrtc.PeerConnectionState) {
	logger := t.GetLogger().WithField("#method", "onSourceConnectionStateChange")

	switch s {
	case webrtc.PeerConnectionStateNew,
		webrtc.PeerConnectionStateConnecting,
		webrtc.PeerConnectionStateConnected,
		webrtc.PeerConnectionStateDisconnected:
	case webrtc.PeerConnectionStateFailed,
		webrtc.PeerConnectionStateClosed:
		go t.Close(t.context())
	}

	logger.WithField("state", s.String()).Tracef("peer connection state changed")
}

func (t *WebrtcTransport) tryNewSysDataChannel() {
	var err error

	logger := t.GetLogger().WithField("#method", "tryNewSysDataChannel")

	t.sysMtx.Lock()
	defer t.sysMtx.Unlock()

	if t.isSysDataChannelAlive() {
		logger.WithField("label", t.sysDataChannel.Label()).Tracef("sys data channel still alive, skip")
		return
	}
	if t.sysDataChannel != nil {
		logger = logger.WithField("old", t.sysDataChannel.Label())
		logger.WithError(t.sysDataChannel.Close()).Debugf("close old sys data channel")
		t.sysDataChannel = nil
		t.rwc = nil
	}

	ordered := !t.enableKcp
	t.sysDataChannel, err = t.pc.CreateDataChannel(t.nextDataChannelLabel("sys"), &webrtc.DataChannelInit{Ordered: &ordered})
	if err != nil {
		defer t.Close(t.context())
		logger.WithError(err).Debugf("failed to new sys data channel")
		return
	}
	logger = logger.WithField("label", t.sysDataChannel.Label())

	t.sysDataChannel.OnOpen(func() { t.onSysDataChannelOpen(t.sysDataChannel) })

	logger.Tracef("new sys data channel")
}

func (t *WebrtcTransport) isSysDataChannelAlive() bool {
	logger := t.GetLogger().WithField("#method", "isSysDataChannelAlive")

	if t.sysDataChannel == nil {
		return false
	}

	ttl, err := t.ping(t.context())
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

		ndc, err := t.pc.CreateDataChannel(label, &webrtc.DataChannelInit{Ordered: &ordered})
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
