package redis_signaling

import (
	"time"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/signaling"
)

func DefaultEngineOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{
		"url":                   "redis://127.0.0.1:6379/0",
		"waitWiredEventTimeout": 13 * time.Second,
		"resolvePeriod":         61 * time.Second,
		"healthCheckPeriod":     57 * time.Second,
	})
}

func WithURL(url string) signaling.NewEngineOption {
	return func(o ofn.Option) {
		o["url"] = url
	}
}

func WithWaitWiredEventTimeout(d time.Duration) signaling.NewEngineOption {
	return func(o ofn.Option) {
		o["waitWiredEventTimeout"] = d
	}
}

func WithResolvePeriod(d time.Duration) signaling.NewEngineOption {
	return func(o ofn.Option) {
		o["resolvePeriod"] = d
	}
}

func WithHealthCheckPeriod(d time.Duration) signaling.NewEngineOption {
	return func(o ofn.Option) {
		o["healthCheckPeriod"] = d
	}
}
