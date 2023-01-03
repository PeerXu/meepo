package listenerer_socks5

import "net"

func (l *Socks5Listener) Serve(lis net.Listener) error {
	logger := l.GetLogger().WithField("#method", "Serve")
	logger.Tracef("socks5 server listening")
	err := l.server.Serve(lis)
	if err != nil {
		logger.WithError(err).Debugf("failed to serve")
		return err
	}
	logger.Tracef("socks5 server terminated")
	return nil
}
