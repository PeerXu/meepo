package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) State() meepo_interface.TransportState {
	if t.isClosed() {
		return meepo_interface.TRANSPORT_STATE_CLOSED
	}

	anyConnecting := false
	anyConnected := false
	t.peerConnections.Range(func(sess Session, pc *webrtc.PeerConnection) bool {
		switch pc.ConnectionState() {
		case webrtc.PeerConnectionStateNew:
		case webrtc.PeerConnectionStateConnecting:
			anyConnecting = true
		case webrtc.PeerConnectionStateConnected:
			anyConnected = true
		case webrtc.PeerConnectionStateDisconnected:
		case webrtc.PeerConnectionStateFailed:
		case webrtc.PeerConnectionStateClosed:
		}
		return true
	})

	if !t.isReady() {
		return meepo_interface.TRANSPORT_STATE_NEW
	}

	if anyConnected {
		return meepo_interface.TRANSPORT_STATE_CONNECTED
	}

	if anyConnecting {
		return meepo_interface.TRANSPORT_STATE_CONNECTING
	}

	return meepo_interface.TRANSPORT_STATE_DISCONNECTED
}
