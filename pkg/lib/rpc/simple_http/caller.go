package rpc_simple_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

type SimpleHttpCaller struct {
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
	logger      logging.Logger
	httpc       *http.Client
	baseURL     string
}

func NewSimpleHttpCaller(opts ...rpc_core.NewCallerOption) (rpc_core.Caller, error) {
	o := option.ApplyWithDefault(DefaultNewSimpleHttpCallerOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	mr, err := marshaler.GetMarshaler(o)
	if err != nil {
		return nil, err
	}

	umr, err := marshaler.GetUnmarshaler(o)
	if err != nil {
		return nil, err
	}

	httpc, err := well_known_option.GetHttpClient(o)
	if err != nil {
		return nil, err
	}

	baseURL, err := GetBaseURL(o)
	if err != nil {
		return nil, err
	}

	return &SimpleHttpCaller{
		marshaler:   mr,
		unmarshaler: umr,
		logger:      logger,
		httpc:       httpc,
		baseURL:     baseURL,
	}, nil
}

func init() {
	rpc_core.RegisterNewCallerFunc("simple_http", NewSimpleHttpCaller)
}
