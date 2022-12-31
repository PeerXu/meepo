package meepo_socks5

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (ss *Socks5Server) GetLogger() logging.Logger {
	return ss.logger.WithFields(logging.Fields{
		"#instance": "Socks5Server",
	})
}
