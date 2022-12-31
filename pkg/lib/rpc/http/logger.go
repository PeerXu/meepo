package rpc_http

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (s *HttpServer) GetLogger() logging.Logger {
	return s.logger.WithFields(logging.Fields{
		"#instance": "HttpServer",
		"listener":  s.listener.Addr().String(),
	})
}

func (c *HttpCaller) GetLogger() logging.Logger {
	return c.logger.WithField("#instance", "HttpCaller")
}
