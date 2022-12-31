package meepo_core

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (mp *Meepo) naviOnce(req *NaviRequest) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":   "naviOnce",
		"tracker":   req.Tracker.String(),
		"candidate": req.Tracker.String(),
	})

	t, err := mp.NewTransport(mp.context(), req.Candidate, WithGetTrackersFunc(func(addr Addr) ([]Tracker, bool, error) {
		tk, err := mp.getTracker(req.Tracker)
		if err != nil {
			return nil, false, err
		}
		return []Tracker{tk}, true, nil
	}))
	if err != nil {
		logger.WithError(err).Debugf("failed to new transport")
		return
	}

	if err = t.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait ready")
		return
	}

	logger.Tracef("navi once")
}
