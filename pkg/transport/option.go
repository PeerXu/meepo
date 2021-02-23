package transport

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
)

type NewTransportOption func(objx.Map)

func WithID(id string) NewTransportOption {
	return func(o objx.Map) {
		o["id"] = id
	}
}

func WithPeerID(id string) NewTransportOption {
	return func(o objx.Map) {
		o["peerID"] = id
	}
}

func WithLogger(logger logrus.FieldLogger) NewTransportOption {
	return func(o objx.Map) {
		o["logger"] = logger
	}
}

type OnDataChannelCreateHandler func(DataChannel)

type CreateDataChannelOption func(objx.Map)

func WithOrdered(ordered bool) CreateDataChannelOption {
	return func(o objx.Map) {
		o["ordered"] = ordered
	}
}
