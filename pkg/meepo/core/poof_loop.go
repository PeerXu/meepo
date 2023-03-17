package meepo_core

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) poofLoop() {
	logger := mp.GetLogger().WithField("#method", "poofLoop")
	defer logger.Tracef("exit poof loop")
	ticker := time.NewTicker(mp.poofInterval)
	for true {
		logger := logger.WithFields(logging.Fields{
			"poofInterval": mp.poofInterval,
			"nextPoofAt":   time.Now().Add(mp.poofInterval),
		})
		ticker.Reset(mp.poofInterval)

		if mp.isClosed() {
			return
		}

		go mp.poofOnce()

		select {
		case <-ticker.C:
			logger.Tracef("poof by ticker tick")
		case <-mp.poofNowCh:
			logger.Tracef("poof by someone kick")
		}
	}
}
