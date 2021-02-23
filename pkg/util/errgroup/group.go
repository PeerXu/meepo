package errgroup

import "sync"

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
