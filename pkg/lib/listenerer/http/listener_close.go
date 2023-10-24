package listenerer_http

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (l *HttpListener) Close() error {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "Close",
	})

	if err := l.lis.Close(); err != nil {
		logger.WithError(err).WithField("listen", l.lis.Addr().String()).Debugf("failed to close listener")
		return err
	}

	l.closeOnce.Do(func() { close(l.conns) })

	logger.Tracef("listener closed")

	return nil
}
