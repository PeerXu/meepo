package rpc_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
)

func DefaultNewHttpCallerOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_HTTP_CLIENT: http.DefaultClient,
	})
}

const (
	OPTION_BASE_URL = "baseURL"
)

var (
	WithBaseURL, GetBaseURL = option.New[string](OPTION_BASE_URL)
)
