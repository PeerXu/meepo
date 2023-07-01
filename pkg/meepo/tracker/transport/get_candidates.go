package tracker_transport

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

func (tk *TransportTracker) GetCandidates(target addr.Addr, requests int, excludes []addr.Addr) (candidates []addr.Addr, err error) {
	ctx := tk.context()
	var excludesStrSlice []string
	for _, x := range excludes {
		excludesStrSlice = append(excludesStrSlice, x.String())
	}
	req := &tracker_interface.GetCandidatesRequest{
		Target:   target.String(),
		Requests: requests,
		Excludes: excludesStrSlice,
	}
	var res tracker_interface.GetCandidatesResponse
	err = tk.transport.Call(ctx, tracker_core.METHOD_GET_CANDIDATES, req, &res)
	if err != nil {
		return nil, err
	}
	for _, candidateStr := range res.Candidates {
		candidate, err := addr.FromString(candidateStr)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return
}
