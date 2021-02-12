package http_api

import (
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/api"
)

func WithHost(host string) api.NewServerOption {
	return func(o objx.Map) {
		o["host"] = host
	}
}

func WithPort(port int32) api.NewServerOption {
	return func(o objx.Map) {
		o["port"] = port
	}
}
