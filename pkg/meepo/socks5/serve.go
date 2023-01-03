package meepo_socks5

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (ss *Socks5Server) Serve(ctx context.Context) <-chan error {
	logger := ss.GetLogger().WithFields(logging.Fields{
		"#method":  "Serve",
		"listener": ss.listener.Addr().String(),
	})
	ss.errors = make(chan error, 1)

	logger.Debugf("serve")
	go func() {
		defer close(ss.errors)
		defer logger.Tracef("terminated")

		if err := ss.server.Serve(ss.listener); err != nil {
			ss.errors <- err
		}
	}()

	return ss.errors
}
