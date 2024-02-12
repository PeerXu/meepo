package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

type (
	GetCandidatesRequest  = tracker_interface.GetCandidatesRequest
	GetCandidatesResponse = tracker_interface.GetCandidatesResponse
)

func (mp *Meepo) onGetCandidates(ctx context.Context, req GetCandidatesRequest) (res GetCandidatesResponse, err error) {
	target, err := addr.FromString(req.Target)
	if err != nil {
		return
	}
	var excludes []Addr
	for _, excludeStr := range req.Excludes {
		var exclude addr.Addr
		exclude, err = addr.FromString(excludeStr)
		if err != nil {
			return
		}
		excludes = append(excludes, exclude)
	}
	candidates, err := mp.getCandidates(target, req.Requests, excludes)
	if err != nil {
		return
	}
	for _, candidate := range candidates {
		res.Candidates = append(res.Candidates, candidate.String())
	}
	return
}

func (mp *Meepo) getCandidates(target Addr, count int, excludes []Addr) (candidates []Addr, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "getCandidates",
		"target":  target.String(),
		"count":   count,
	})

	tks, _, err := mp.getCloserTrackers(target, count, excludes)
	if err != nil {
		logger.WithError(err).Debugf("failed to get closer trackers")
		return nil, err
	}

	for _, tk := range tks {
		candidates = append(candidates, tk.Addr())
	}

	logger.Tracef("on get candidates")

	return candidates, nil
}
