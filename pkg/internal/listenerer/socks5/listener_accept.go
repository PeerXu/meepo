package listenerer_socks5

import (
	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (l *Socks5Listener) Accept() (listenerer_interface.Conn, error) {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "Accept",
	})

	conn, ok := <-l.conns
	if !ok {
		err := ErrUseClosedNetworkConnection
		logger.WithError(err).Debugf("accept from closed network connection")
		return nil, err
	}

	logger.WithFields(logging.Fields{
		"localNetwork":  conn.LocalAddr().Network(),
		"localAddress":  conn.LocalAddr().String(),
		"remoteNetwork": conn.RemoteAddr().Network(),
		"remoteAddress": conn.RemoteAddr().String(),
	}).Tracef("accept")

	return conn, nil
}
