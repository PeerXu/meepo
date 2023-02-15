package tracker_interface

import "github.com/pion/webrtc/v3"

type AddPeerConnectionRequest struct {
	Session int32
	Offer   webrtc.SessionDescription
}

type AddPeerConnectionResponse struct {
	Session int32
	Answer  webrtc.SessionDescription
}
