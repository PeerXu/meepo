package rpc_simple_http

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (c *SimpleHttpCaller) GetLogger() logging.Logger {
	return c.logger.WithFields(logging.Fields{"#instance": "SimpleHttpCaller"})
}

func (s *SimpleHttpServer) GetLogger() logging.Logger {
	return s.logger.WithFields(logging.Fields{"#instance": "SimpleHttpServer"})
}
