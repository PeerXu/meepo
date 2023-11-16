package listenerer_http

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (l *HttpListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := l.GetLogger().WithFields(logging.Fields{
		"#method": "ServeHTTP",
	})

	switch r.Method {
	case http.MethodGet:
		l.handleGet(w, r)
	case http.MethodConnect:
		l.handleConnect(w, r)
	default:
		logger.WithField("method", r.Method).Warningf("method not allowed")
		l.writeStatusCode(w, http.StatusMethodNotAllowed)
	}
}
