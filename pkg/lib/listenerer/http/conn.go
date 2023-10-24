package listenerer_http

import (
	"io"
	"net/http"
	"sync"
)

type HttpConn struct {
	reader    io.Reader
	writer    io.Writer
	request   *http.Request
	close     chan struct{}
	closeOnce sync.Once
}
