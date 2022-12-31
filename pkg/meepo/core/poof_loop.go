package meepo_core

import "time"

func (mp *Meepo) poofLoop() {
	logger := mp.GetLogger().WithField("#method", "poofLoop")
	defer logger.Tracef("exit poof loop")
	ticker := time.NewTicker(mp.poofInterval)
	for ; true; <-ticker.C {
		if mp.isClosed() {
			return
		}
		go mp.poofOnce()
	}
}
