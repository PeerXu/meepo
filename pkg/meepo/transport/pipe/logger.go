package transport_pipe

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (t *PipeTransport) GetLogger() logging.Logger {
	return t.logger.WithFields(logging.Fields{
		"#instance": "PipeTransport",
		"addr":      t.Addr(),
		"session":   t.Session(),
	})
}

func (t *PipeTransport) GetRawLogger() logging.Logger {
	return t.logger
}

func (c *PipeChannel) GetLogger() logging.Logger {
	return c.logger.WithFields(logging.Fields{
		"#instance": "PipeChannel",
		"channelID": c.ID(),
	})
}
