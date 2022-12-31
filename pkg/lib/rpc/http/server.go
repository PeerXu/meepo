package rpc_http

import (
	"net"
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

type HttpServer struct {
	handler     rpc_core.Handler
	signer      crypto_interface.Signer
	cryptor     crypto_interface.Cryptor
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
	listener    net.Listener
	httpd       *http.Server
	errors      chan error
	logger      logging.Logger
}

func NewHttpServer(opts ...rpc_core.NewServerOption) (rpc_core.Server, error) {
	o := option.Apply(opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	handler, err := rpc_core.GetHandler(o)
	if err != nil {
		return nil, err
	}

	cryptor, err := crypto_core.GetCryptor(o)
	if err != nil {
		return nil, err
	}

	signer, err := crypto_core.GetSigner(o)
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

	listener, err := well_known_option.GetListener(o)
	if err != nil {
		return nil, err
	}

	httpd := &http.Server{}

	s := &HttpServer{
		handler:     handler,
		signer:      signer,
		cryptor:     cryptor,
		marshaler:   mr,
		unmarshaler: umr,
		listener:    listener,
		httpd:       httpd,
		logger:      logger,
	}
	httpd.Handler = s.Routers()

	return s, nil
}

func init() {
	rpc_core.RegisterNewServerFunc("http", NewHttpServer)
}
