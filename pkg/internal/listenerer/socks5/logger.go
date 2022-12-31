package listenerer_socks5

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (l *Socks5Listener) GetLogger() logging.Logger {
	return l.logger.WithFields(logging.Fields{
		"#instance": "Socks5Listener",
		"addr":      l.Addr().String(),
	})
}
