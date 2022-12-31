package transport_webrtc

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (t *WebrtcTransport) GetLogger() logging.Logger {
	return t.logger.WithFields(logging.Fields{
		"#instance": "WebrtcTransport",
		"addr":      t.Addr(),
		"role":      t.role,
		"enableMux": t.enableMux,
		"enableKcp": t.enableKcp,
	})
}

func (t *WebrtcTransport) GetRawLogger() logging.Logger {
	return t.logger
}

func (t *WebrtcTransport) wrapMessage(m Message) logging.Fields {
	return logging.Fields{
		"session": m.Session,
		"scope":   m.Scope,
		"method":  m.Method,
	}
}

func (c *WebrtcSourceChannel) GetLogger() logging.Logger {
	return c.logger.WithFields(logging.Fields{
		"#instance": "WebrtcSourceChannel",
		"channelID": c.ID(),
	})
}

func (c *WebrtcSinkChannel) GetLogger() logging.Logger {
	return c.logger.WithFields(logging.Fields{
		"#instance": "WebrtcSinkChannel",
		"channelID": c.ID(),
	})
}
