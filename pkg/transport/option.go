package transport

import (
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/ofn"
)

type NewTransportOption = ofn.OFN

func WithID(id string) NewTransportOption {
	return func(o ofn.Option) {
		o["id"] = id
	}
}

func WithPeerID(id string) NewTransportOption {
	return func(o ofn.Option) {
		o["peerID"] = id
	}
}

func WithLogger(logger logrus.FieldLogger) NewTransportOption {
	return func(o ofn.Option) {
		o["logger"] = logger
	}
}

type OnDataChannelCreateHandler func(DataChannel)

type CreateDataChannelOption = ofn.OFN

func WithOrdered(ordered bool) CreateDataChannelOption {
	return func(o ofn.Option) {
		o["ordered"] = ordered
	}
}
