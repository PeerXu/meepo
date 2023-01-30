package meepo_core

import (
	"math/rand"

	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_routing_table_interface "github.com/PeerXu/meepo/pkg/meepo/routing_table/interface"
)

func (mp *Meepo) poofOnce() {
	logger := mp.GetLogger().WithField("#method", "poofOnce")
	if !mp.poofMtx.TryLock() {
		logger.Tracef("other goroutine is poofing")
		return
	}
	defer mp.poofMtx.Unlock()

	hr := mp.routingTable.HealthReport()
	defer mp.briefHealthReport(logger, hr)

	var lvl meepo_routing_table_interface.HealthLevel
	var buckets []int
	if hr.Summary[meepo_routing_table_interface.HEALTH_LEVEL_RED] > 0 {
		buckets = hr.Report[meepo_routing_table_interface.HEALTH_LEVEL_RED]
		lvl = meepo_routing_table_interface.HEALTH_LEVEL_RED
	} else if hr.Summary[meepo_routing_table_interface.HEALTH_LEVEL_YELLOW] > 0 {
		buckets = hr.Report[meepo_routing_table_interface.HEALTH_LEVEL_YELLOW]
		lvl = meepo_routing_table_interface.HEALTH_LEVEL_YELLOW
	} else {
		return
	}
	cpl := buckets[rand.Intn(len(buckets)-1)]
	targetID := mp.routingTable.GenRandIDWithCpl(cpl)
	targetAddr := ID2Addr(targetID)
	logger = logger.WithFields(logging.Fields{
		"bucket": cpl,
		"target": targetAddr.String(),
		"level":  lvl.String(),
	})

	tks, _, err := mp.getCloserTrackers(targetAddr, mp.poofCount, nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to get closer trackers")
		return
	}

	var handledCandidates []Addr
	mtx := lock.NewLock(well_known_option.WithName("poofOnceMtx"))
	for _, tk := range tks {
		go func(tk Tracker) {
			logger := logger.WithField("tracker", tk.Addr().String())
			candidates, err := tk.GetCandidates(targetAddr, mp.poofCount, []Addr{mp.Addr()})
			if err != nil {
				logger.WithError(err).Debugf("failed to get candidates")
				return
			}
			if !mp.existTransport(tk.Addr()) {
				candidates = append(candidates, tk.Addr())
			}

			mtx.Lock()
			defer mtx.Unlock()
			for _, candidate := range candidates {
				if mp.isClosed() {
					return
				}

				if ContainAddr(handledCandidates, candidate) ||
					mp.Addr().Equal(candidate) ||
					mp.existTransport(candidate) {
					continue
				}

				handledCandidates = append(handledCandidates, candidate)
				mp.naviRequests <- &NaviRequest{
					Candidate: candidate,
					Tracker:   tk.Addr(),
				}
				logger.WithField("candidate", candidate.String()).Tracef("create navi request")
			}
		}(tk)
	}
}

func (mp *Meepo) briefHealthReport(logger logging.Logger, hr *meepo_routing_table_interface.HealthReport) {
	logger.WithFields(logging.Fields{
		meepo_routing_table_interface.HEALTH_LEVEL_RED.String():    hr.Summary[meepo_routing_table_interface.HEALTH_LEVEL_RED],
		meepo_routing_table_interface.HEALTH_LEVEL_YELLOW.String(): hr.Summary[meepo_routing_table_interface.HEALTH_LEVEL_YELLOW],
		meepo_routing_table_interface.HEALTH_LEVEL_GREEN.String():  hr.Summary[meepo_routing_table_interface.HEALTH_LEVEL_GREEN],
	}).Tracef("health report")
}
