package rpc_simple_http

import (
	"net"
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

type SimpleHttpServer struct {
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
	handler     rpc_core.Handler
	listener    net.Listener
	httpd       *http.Server
	errors      chan error
	logger      logging.Logger
}

func NewSimpleHttpServer(opts ...rpc_core.NewServerOption) (rpc_core.Server, error) {
	o := option.Apply(opts...)

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

	handler, err := rpc_core.GetHandler(o)
	if err != nil {
		return nil, err
	}

	listener, err := well_known_option.GetListener(o)
	if err != nil {
		return nil, err
	}

	httpd := &http.Server{}

	s := &SimpleHttpServer{
		marshaler:   mr,
		unmarshaler: umr,
		handler:     handler,
		listener:    listener,
		httpd:       httpd,
		logger:      logger,
	}
	httpd.Handler = s.Routers()

	return s, nil
}

func init() {
	rpc_core.RegisterNewServerFunc("simple_http", NewSimpleHttpServer)
}
