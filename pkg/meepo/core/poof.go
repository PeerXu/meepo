package meepo_core

import "time"

func (mp *Meepo) getPoofInterval() time.Duration {
	return mp.poofInterval
}
