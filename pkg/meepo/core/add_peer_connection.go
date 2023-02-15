package meepo_core

import (
	"github.com/pion/webrtc/v3"

	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

func (mp *Meepo) newAddPeerConnectionRequest(target Addr, sess transport_webrtc.Session, offer webrtc.SessionDescription) (*crypto_core.Packet, error) {
	req := &tracker_interface.AddPeerConnectionRequest{
		Session: int32(sess),
		Offer:   offer,
	}
	return mp.encryptMessage(target, req)
}
