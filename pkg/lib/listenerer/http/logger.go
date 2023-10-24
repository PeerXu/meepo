package listenerer_http

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (l *HttpListener) GetLogger() logging.Logger {
	return l.logger.WithFields(logging.Fields{
		"#instance": "HttpListener",
		"addr":      l.Addr().String(),
	})
}
