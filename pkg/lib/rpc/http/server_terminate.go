package rpc_http

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (s *HttpServer) Terminate(ctx context.Context) error {
	logger := s.GetLogger().WithFields(logging.Fields{
		"#method": "Terminate",
	})
	logger.Tracef("terminating")
	return s.httpd.Shutdown(ctx)
}
