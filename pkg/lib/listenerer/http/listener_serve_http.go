package listenerer_http

import (
	"fmt"
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (l *HttpListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "onConnect",
	})

	if r.Method != http.MethodConnect {
		l.writeStatusCode(w, http.StatusMethodNotAllowed)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		l.writeStatusCode(w, http.StatusNotAcceptable)
		return
	}

	hijacked, _, err := hijacker.Hijack()
	if err != nil {
		l.writeStatusCode(w, http.StatusInternalServerError)
		return
	}

	conn := &HttpConn{
		reader:  hijacked,
		writer:  hijacked,
		request: r,
		close:   make(chan struct{}),
	}
	defer close(conn.close)

	l.conns <- conn

	if err = conn.WaitEnabled(l.connWaitEnabledTimeout); err != nil {
		logger.WithError(err).Debugf("failed to wait conn enabled")
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%s 200 Connection established\r\n\r\n", r.Proto)))
	if err != nil {
		logger.WithError(err).Debugf("failed to reply to http client")
		return
	}

	<-conn.close

	logger.Tracef("remote http connection closed")

	return
}
