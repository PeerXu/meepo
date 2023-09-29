package rpc_simple_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

type ctxSession string

const (
	OPTION_BASE_URL = "baseURL"
)

var (
	WithBaseURL, GetBaseURL = option.New[string](OPTION_BASE_URL)
)

func DefaultNewSimpleHttpCallerOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_HTTP_CLIENT: http.DefaultClient,
	})
}
