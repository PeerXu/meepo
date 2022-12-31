package meepo_core

import (
	"github.com/pion/webrtc/v3"
)

func (mp *Meepo) newPeerConnection() (*webrtc.PeerConnection, error) {
	return mp.webrtcAPI.NewPeerConnection(mp.webrtcConfiguration)
}
