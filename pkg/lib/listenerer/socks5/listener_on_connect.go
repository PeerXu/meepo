package listenerer_socks5

import (
	"context"
	"io"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (l *Socks5Listener) onConnect(ctx context.Context, writer io.Writer, request *socks5.Request) (err error) {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "onConnect",
	})

	defer func() {
		var ok bool
		if err, ok = recover().(error); ok {
			logger.WithError(err).Debugf("recover from panic")
			if er := socks5.SendReply(writer, statute.RepServerFailure, nil); er != nil {
				logger.WithError(er).Debugf("faile to send reply to socks5 client")
			}
		}
	}()

	conn := &Socks5Conn{
		writer:  writer,
		request: request,
		close:   make(chan struct{}),
	}
	defer close(conn.close)

	l.conns <- conn

	if err = conn.WaitEnabled(l.connWaitEnabledTimeout); err != nil {
		logger.WithError(err).Debugf("failed to wait conn enabled")
		return err
	}

	if err = socks5.SendReply(writer, statute.RepSuccess, request.LocalAddr); err != nil {
		logger.WithError(err).Debugf("failed to send success reply to socks5 client")
		return err
	}

	<-conn.close

	logger.Tracef("remote socks5 conenction closed")

	return nil
}
