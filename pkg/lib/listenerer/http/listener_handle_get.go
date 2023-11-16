package listenerer_http

import (
	"io"
	"net"
	"net/http"
	"strings"
)

func (l *HttpListener) handleGet(w http.ResponseWriter, r *http.Request) {
	or := new(http.Request)
	*or = *r

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if proir, ok := or.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(proir, ", ") + ", " + clientIP
		}
		or.Header.Set("X-Forwarded-For", clientIP)
	}

	ow, err := l.transport.RoundTrip(or)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer ow.Body.Close()

	for key, vals := range ow.Header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(ow.StatusCode)
	io.Copy(w, ow.Body)
}
