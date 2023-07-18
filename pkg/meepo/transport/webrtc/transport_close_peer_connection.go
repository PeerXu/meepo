//go:build !(js && wasm)

package transport_webrtc

import "github.com/pion/webrtc/v3"

func (t *WebrtcTransport) closePeerConnection(sess Session, pc *webrtc.PeerConnection) error {
	return pc.Close()
}
