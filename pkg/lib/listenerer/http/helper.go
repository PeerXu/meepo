package listenerer_http

import "net/http"

func (l *HttpListener) writeStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}
