package redis_signaling

import (
	"time"

	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
)

func DefaultEngineOption() objx.Map {
	return objx.New(map[string]interface{}{
		"url":                   "redis://127.0.0.1:6379/0",
		"waitWiredEventTimeout": 13 * time.Second,
		"resolvePeriod":         61 * time.Second,
		"healthCheckPeriod":     57 * time.Second,
	})
}

func WithURL(url string) signaling.NewEngineOption {
	return func(o objx.Map) {
		o["url"] = url
	}
}

func WithWaitWiredEventTimeout(d time.Duration) signaling.NewEngineOption {
	return func(o objx.Map) {
		o["waitWiredEventTimeout"] = d
	}
}

func WithResolvePeriod(d time.Duration) signaling.NewEngineOption {
	return func(o objx.Map) {
		o["resolvePeriod"] = d
	}
}

func WithHealthCheckPeriod(d time.Duration) signaling.NewEngineOption {
	return func(o objx.Map) {
		o["healthCheckPeriod"] = d
	}
}
