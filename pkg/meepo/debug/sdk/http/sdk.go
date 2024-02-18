package meepo_debug_sdk_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_debug_interface "github.com/PeerXu/meepo/pkg/meepo/debug/interface"
	meepo_debug_sdk_core "github.com/PeerXu/meepo/pkg/meepo/debug/sdk/core"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
}

func NewSDK(opts ...HttpSDKOption) (meepo_debug_interface.MeepoDebugInterface, error) {
	o := option.ApplyWithDefault(defaultHttpSDKOption(), opts...)

	baseUrl, err := GetBaseURL(o)
	if err != nil {
		return nil, err
	}

	httpClient, err := GetHTTPClient(o)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}, nil
}

func init() {
	meepo_debug_sdk_core.Register("http", NewSDK)
}
