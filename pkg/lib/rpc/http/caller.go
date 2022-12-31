package rpc_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

type HttpCaller struct {
	signer      crypto_interface.Signer
	cryptor     crypto_interface.Cryptor
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler

	logger  logging.Logger
	httpc   *http.Client
	baseURL string
}

func NewHttpCaller(opts ...rpc_core.NewCallerOption) (rpc_core.Caller, error) {
	o := option.ApplyWithDefault(DefaultNewHttpCallerOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	signer, err := crypto_core.GetSigner(o)
	if err != nil {
		return nil, err
	}

	cryptor, err := crypto_core.GetCryptor(o)
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

	client, err := well_known_option.GetHttpClient(o)
	if err != nil {
		return nil, err
	}

	baseURL, err := GetBaseURL(o)
	if err != nil {
		return nil, err
	}

	return &HttpCaller{
		signer:      signer,
		cryptor:     cryptor,
		marshaler:   mr,
		unmarshaler: umr,
		logger:      logger,
		httpc:       client,
		baseURL:     baseURL,
	}, nil
}

func init() {
	rpc_core.RegisterNewCallerFunc("http", NewHttpCaller)
}
