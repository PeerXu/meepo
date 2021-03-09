package group

import (
	"sync"
)

func DONE() (interface{}, error) { return nil, nil }

type Group interface {
	Go(func() (interface{}, error), ...func(interface{}, error))
	Wait() (interface{}, error)
}

type NewGroupFunc func() Group

var (
	NewAllGroupFunc  = func() Group { return new(AllGroup) }
	NewRaceGroupFunc = func() Group { return new(RaceGroup) }
	NewAnyGroupFunc  = func() Group { return new(AnyGroup) }
)

type AllGroup struct {
	initOnce sync.Once

	doneWg   sync.WaitGroup
	fastDone chan struct{}

	errOnce sync.Once
	err     error
}

func (g *AllGroup) Go(f func() (interface{}, error), cb ...func(interface{}, error)) {
	g.ensureInit()

	g.doneWg.Add(1)
	go func() {
		defer g.doneWg.Done()

		if _, err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if len(cb) > 0 {
					cb[0](nil, err)
				}
				close(g.fastDone)
			})
		}
	}()
}

func (g *AllGroup) Wait() (interface{}, error) {
	g.ensureInit()
	defer g.ensureRelease()

	allDone := make(chan struct{})
	go func() {
		g.doneWg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
	case <-g.fastDone:
	}

	return nil, g.err
}

func (g *AllGroup) ensureInit() {
	g.initOnce.Do(func() {
		g.fastDone = make(chan struct{})
	})
}

func (g *AllGroup) ensureRelease() {
	g.errOnce.Do(func() { close(g.fastDone) })
}

type RaceGroup struct {
	initOnce sync.Once

	done chan struct{}

	doneOnce sync.Once
	ret      interface{}
	err      error
}

func (g *RaceGroup) ensureInit() {
	g.initOnce.Do(func() {
		g.done = make(chan struct{})
	})

}

func (g *RaceGroup) Go(f func() (interface{}, error), cb ...func(interface{}, error)) {
	g.ensureInit()

	go func() {
		ret, err := f()
		g.doneOnce.Do(func() {
			g.ret = ret
			g.err = err
			if len(cb) > 0 {
				cb[0](ret, err)
			}
			close(g.done)
		})
	}()
}

func (g *RaceGroup) Wait() (interface{}, error) {
	g.ensureInit()

	select {
	case <-g.done:
	}

	return g.ret, g.err
}

type AnyGroup struct {
	initOnce sync.Once

	doneWg   sync.WaitGroup
	fastDone chan struct{}

	errOnce sync.Once
	err     error

	retOnce sync.Once
	ret     interface{}
}

func (g *AnyGroup) ensureInit() {
	g.initOnce.Do(func() {
		g.fastDone = make(chan struct{})
	})
}

func (g *AnyGroup) ensureRelease() {
	g.retOnce.Do(func() { close(g.fastDone) })
}

func (g *AnyGroup) Wait() (interface{}, error) {
	g.ensureInit()
	defer g.ensureRelease()

	allDone := make(chan struct{})
	go func() {
		g.doneWg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
		return g.ret, g.err
	case <-g.fastDone:
		return g.ret, nil
	}
}

func (g *AnyGroup) Go(f func() (interface{}, error), cb ...func(interface{}, error)) {
	g.ensureInit()

	g.doneWg.Add(1)
	go func() {
		defer g.doneWg.Done()

		ret, err := f()
		if err != nil {
			g.errOnce.Do(func() { g.err = err })
			return
		}

		g.retOnce.Do(func() {
			g.ret = ret
			if len(cb) > 0 {
				cb[0](ret, nil)
			}
			close(g.fastDone)
		})
	}()
}
