//go:build js && wasm

package transport_webrtc

import "github.com/pion/webrtc/v3"

func (t *WebrtcTransport) closePeerConnection(sess Session, pc *webrtc.PeerConnection) error {
	err := pc.Close()
	if err != nil {
		return err
	}

	t.unregisterPeerConnection(sess)

	return nil
}
