package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) State() meepo_interface.TransportState {
	if t.isClosed() {
		return meepo_interface.TRANSPORT_STATE_CLOSED
	}

	connected := false
	t.peerConnections.Range(func(sess Session, pc *webrtc.PeerConnection) bool {
		switch pc.ConnectionState() {
		case webrtc.PeerConnectionStateNew:
		case webrtc.PeerConnectionStateConnecting:
		case webrtc.PeerConnectionStateConnected:
			connected = true
			return false
		case webrtc.PeerConnectionStateDisconnected:
		case webrtc.PeerConnectionStateFailed:
		case webrtc.PeerConnectionStateClosed:
		}
		return true
	})

	if connected {
		return meepo_interface.TRANSPORT_STATE_CONNECTED
	}

	if !t.isReady() {
		return meepo_interface.TRANSPORT_STATE_NEW
	}

	if !t.isConnectedOnce() {
		return meepo_interface.TRANSPORT_STATE_CONNECTING
	}

	return meepo_interface.TRANSPORT_STATE_DISCONNECTED
}

func (t *WebrtcTransport) isConnectedOnce() bool {
	return t.connectedOnce.Load()
}
