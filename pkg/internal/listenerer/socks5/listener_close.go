package listenerer_socks5

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (l *Socks5Listener) Close() error {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "Close",
	})

	l.closeOnce.Do(func() { close(l.conns) })

	if err := l.lis.Close(); err != nil {
		logger.WithError(err).WithField("listen", l.lis.Addr().String()).Debugf("failed to close listener")
		return err
	}

	logger.Tracef("listener closed")

	return nil
}
