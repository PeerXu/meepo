package meepo_core

import "time"

func (mp *Meepo) increasePoofInterval() {
	pi := time.Duration(float64(mp.poofInterval) * mp.poofIntervalFactor)
	if pi > mp.poofMaxInterval {
		pi = mp.poofMaxInterval
	}
	mp.poofInterval = pi
}

func (mp *Meepo) decreasePoofInterval() {
	pi := time.Duration(float64(mp.poofInterval) / mp.poofIntervalFactor)
	if pi < mp.poofMinInterval {
		pi = mp.poofMinInterval
	}
	mp.poofInterval = pi
}

func (mp *Meepo) resetPoofInterval() {
	mp.poofInterval = mp.poofMinInterval
}

func (mp *Meepo) poofNow() {
	panic("unimplemented")
}
