package listenerer_socks5

import (
	"io"
	"sync"

	"github.com/things-go/go-socks5"
)

type Socks5Conn struct {
	writer     io.Writer
	request    *socks5.Request
	close      chan struct{}
	closeOnce  sync.Once
	enable     chan struct{}
	enableOnce sync.Once
}
