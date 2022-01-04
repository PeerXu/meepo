package http_api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"

	"github.com/PeerXu/meepo/pkg/api"
	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/ofn"
)

func NewHttpServerOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{
		"host": "127.0.0.1",
		"port": 12345,
	})
}

type HttpServer struct {
	opt ofn.Option

	httpd *http.Server
	eg    errgroup.Group

	meepo *meepo.Meepo
}

func (s *HttpServer) Start(ctx context.Context) error {
	host := cast.ToString(s.opt.Get("host").Inter())
	port := cast.ToString(s.opt.Get("port").Inter())

	s.httpd = &http.Server{
		Addr: net.JoinHostPort(host, port),
	}
	s.httpd.Handler = s.getRouter()

	s.eg.Go(s.httpd.ListenAndServe)

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.httpd.Shutdown(ctx)
}

func (s *HttpServer) Wait() error {
	return s.eg.Wait()
}

func NewHttpServer(opts ...api.NewServerOption) (api.Server, error) {
	o := NewHttpServerOption()
	for _, opt := range opts {
		opt(o)
	}

	mp, ok := o.Get("meepo").Inter().(*meepo.Meepo)
	if !ok {
		return nil, fmt.Errorf("require meepo")
	}

	return &HttpServer{
		opt:   o,
		meepo: mp,
	}, nil
}

func init() {
	api.RegisterNewServerFunc("http", NewHttpServer)
}
