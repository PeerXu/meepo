package signaling

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/ofn"
)

type NewEngineOption = ofn.OFN

func WithID(id string) NewEngineOption {
	return func(o objx.Map) {
		o["id"] = id
	}
}

func WithLogger(logger logrus.FieldLogger) NewEngineOption {
	return func(o objx.Map) {
		o["logger"] = logger
	}
}
