package signaling

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
)

type NewEngineOption func(objx.Map)

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
