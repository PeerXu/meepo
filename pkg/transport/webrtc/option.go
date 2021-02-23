package webrtc_transport

import (
	"github.com/pion/webrtc/v3"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/transport"
)

func AsAnswerer() transport.NewTransportOption {
	return func(o objx.Map) {
		o["role"] = "answerer"
	}
}

func AsOfferer() transport.NewTransportOption {
	return func(o objx.Map) {
		o["role"] = "offerer"
	}
}

func WithWebrtcAPI(api *webrtc.API) transport.NewTransportOption {
	return func(o objx.Map) {
		o["webrtcAPI"] = api
	}
}

func WithICEServers(iceServers []string) transport.NewTransportOption {
	return func(o objx.Map) {
		o["iceServers"] = iceServers
	}
}

func WithOffer(offer *webrtc.SessionDescription) transport.NewTransportOption {
	return func(o objx.Map) {
		o["offer"] = offer
	}
}

func WithOfferHook(offerHook OfferHook) transport.NewTransportOption {
	return func(o objx.Map) {
		o["offerHook"] = offerHook
	}
}

func WithAnswerHook(answerHook AnswerHook) transport.NewTransportOption {
	return func(o objx.Map) {
		o["answerHook"] = answerHook
	}
}
