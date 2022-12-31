package meepo_socks5

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (ss *Socks5Server) Terminate(ctx context.Context) error {
	logger := ss.GetLogger().WithFields(logging.Fields{
		"#method": "Terminate",
	})
	logger.Tracef("terminating")
	return ss.listener.Close()
}
