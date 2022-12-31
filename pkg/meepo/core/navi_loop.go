package meepo_core

func (mp *Meepo) naviLoop() {
	logger := mp.GetLogger().WithField("#method", "naviLoop")
	defer logger.Tracef("naviLoop closed")

	for {
		req, ok := <-mp.naviRequests
		if !ok {
			return
		}

		go mp.naviOnce(req)
	}
}
