package meepo_core

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) naviLoop() {
	logger := mp.GetLogger().WithField("#method", "naviLoop")
	defer logger.Tracef("naviLoop closed")

	for {
		req, ok := <-mp.naviRequests
		if !ok {
			return
		}

		queued := time.Since(req.CreatedAt)
		if queued > mp.naviRequestQueueTimeout {
			logger.WithFields(logging.Fields{
				"requestSession": req.Session,
				"createdAt":      req.CreatedAt,
				"queued":         queued,
			}).Debugf("navi request timeout")
		}

		go mp.naviOnce(req)
	}
}
