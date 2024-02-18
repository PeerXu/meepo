package meepo_debug_sdk_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_BASE_URL    = "baseURL"
	OPTION_HTTP_CLIENT = "httpClient"
)

type HttpSDKOption = option.ApplyOption

func defaultHttpSDKOption() option.Option {
	return option.Option{
		OPTION_HTTP_CLIENT: http.DefaultClient,
	}
}

var (
	WithBaseURL, GetBaseURL       = option.New[string](OPTION_BASE_URL)
	WithHTTPClient, GetHTTPClient = option.New[*http.Client](OPTION_HTTP_CLIENT)
)
