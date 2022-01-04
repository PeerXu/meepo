package http_api

import (
	"github.com/PeerXu/meepo/pkg/api"
	"github.com/PeerXu/meepo/pkg/ofn"
)

func WithHost(host string) api.NewServerOption {
	return func(o ofn.Option) {
		o["host"] = host
	}
}

func WithPort(port int32) api.NewServerOption {
	return func(o ofn.Option) {
		o["port"] = port
	}
}
