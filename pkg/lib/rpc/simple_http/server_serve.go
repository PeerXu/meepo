package rpc_simple_http

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (s *SimpleHttpServer) Serve(context.Context) <-chan error {
	logger := s.GetLogger().WithFields(logging.Fields{
		"#method":  "Serve",
		"listener": s.listener.Addr().String(),
	})
	s.errors = make(chan error, 1)

	logger.Infof("serve")
	go func() {
		defer close(s.errors)
		defer logger.Tracef("terminated")

		if err := s.httpd.Serve(s.listener); err != nil {
			s.errors <- err
		}
	}()

	return s.errors
}
