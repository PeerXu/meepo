package rpc_http

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (s *HttpServer) Serve(ctx context.Context) <-chan error {
	logger := s.GetLogger().WithFields(logging.Fields{
		"#method":  "Serve",
		"listener": s.listener.Addr().String(),
	})
	s.errors = make(chan error, 1)

	logger.Debugf("serve")
	go func() {
		defer close(s.errors)
		defer logger.Tracef("terminated")

		if err := s.httpd.Serve(s.listener); err != nil {
			s.errors <- err
		}
	}()

	return s.errors
}
