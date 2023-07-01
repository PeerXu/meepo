package meepo_core

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) poofLoop() {
	logger := mp.GetLogger().WithField("#method", "poofLoop")
	defer logger.Tracef("exit")
	timer := time.NewTicker(mp.getPoofInterval())
	for {
		pi := mp.getPoofInterval()
		logger := logger.WithFields(logging.Fields{
			"poofInterval": pi,
			"nextPoofAt":   time.Now().Add(pi),
		})
		timer.Reset(pi)

		if mp.isClosed() {
			return
		}

		go mp.poofOnce()

		select {
		case <-timer.C:
			logger.Tracef("poof by ticker timer")
		}
	}
}
