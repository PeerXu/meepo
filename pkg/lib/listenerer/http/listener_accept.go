package listenerer_http

import (
	listenerer_core "github.com/PeerXu/meepo/pkg/lib/listenerer/core"
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (l *HttpListener) Accept() (listenerer_interface.Conn, error) {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "Accept",
	})

	conn, ok := <-l.conns
	if !ok {
		err := listenerer_core.ErrListenerClosed
		logger.WithError(err).Debugf("accept from closed listener")
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
