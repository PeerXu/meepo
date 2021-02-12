package meepo

import (
	"net"
	"sync"
)

func isListenableAddr(addr net.Addr) (net.Addr, bool) {
	lis, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		return nil, false
	}
	defer lis.Close()

	tcpAddr, _ := net.ResolveTCPAddr(lis.Addr().Network(), lis.Addr().String())

	return tcpAddr, true
}

// TODO(Peer): More robustness
func getListenableAddr() net.Addr {
	for {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
		if addr, ok := isListenableAddr(addr); ok {
			return addr
		}
	}
}

// Mutated from errgroup.Group
type ImmediatelyErrorGroup struct {
	initDoneOnce sync.Once
	doneOnce     sync.Once
	done         chan struct{}

	errOnce sync.Once
	err     error
}

func (g *ImmediatelyErrorGroup) Wait() error {
	<-g.done
	return g.err
}

func (g *ImmediatelyErrorGroup) Go(f func() error) {
	g.initDoneOnce.Do(func() {
		g.done = make(chan struct{}, 1)
	})

	go func() {
		defer g.doneOnce.Do(func() { close(g.done) })

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}
