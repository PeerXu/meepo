package webrtc_transport

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/transport"
)

func AsAnswerer() transport.NewTransportOption {
	return func(o ofn.Option) {
		o["role"] = "answerer"
	}
}

func AsOfferer() transport.NewTransportOption {
	return func(o ofn.Option) {
		o["role"] = "offerer"
	}
}

func WithWebrtcAPI(api *webrtc.API) transport.NewTransportOption {
	return func(o ofn.Option) {
		o["webrtcAPI"] = api
	}
}

func WithICEServers(iceServers []string) transport.NewTransportOption {
	return func(o ofn.Option) {
		o["iceServers"] = iceServers
	}
}

func WithOffer(offer *webrtc.SessionDescription) transport.NewTransportOption {
	return func(o ofn.Option) {
		o["offer"] = offer
	}
}

func WithOfferHook(offerHook OfferHook) transport.NewTransportOption {
	return func(o ofn.Option) {
		o["offerHook"] = offerHook
	}
}

func WithAnswerHook(answerHook AnswerHook) transport.NewTransportOption {
	return func(o ofn.Option) {
		o["answerHook"] = answerHook
	}
}
