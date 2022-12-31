package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

type (
	GetCandidatesRequest  = tracker_interface.GetCandidatesRequest
	GetCandidatesResponse = tracker_interface.GetCandidatesResponse
)

func (mp *Meepo) newOnGetCandidatesRequest() any { return &GetCandidatesRequest{} }

func (mp *Meepo) hdrOnGetCandidates(ctx context.Context, _req any) (any, error) {
	req := _req.(*GetCandidatesRequest)
	target, err := addr.FromString(req.Target)
	if err != nil {
		return nil, err
	}
	var excludes []Addr
	for _, excludeStr := range req.Excludes {
		exclude, err := addr.FromString(excludeStr)
		if err != nil {
			return nil, err
		}
		excludes = append(excludes, exclude)
	}
	candidates, err := mp.onGetCandidates(target, req.Count, excludes)
	if err != nil {
		return nil, err
	}
	var res GetCandidatesResponse
	for _, candidate := range candidates {
		res.Candidates = append(res.Candidates, candidate.String())
	}
	return &res, nil
}

func (mp *Meepo) onGetCandidates(target Addr, count int, excludes []Addr) (candidates []Addr, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onGetCandidates",
		"target":  target.String(),
		"count":   count,
	})

	tks, _, err := mp.getNearestTrackers(target, count, excludes)
	if err != nil {
		logger.WithError(err).Debugf("failed to get nearest trackers")
		return nil, err
	}

	for _, tk := range tks {
		candidates = append(candidates, tk.Addr())
	}

	logger.Tracef("on get candidates")

	return candidates, nil
}
