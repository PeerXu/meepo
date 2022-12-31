package tracker_transport

import (
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

func (tk *TransportTracker) GetCandidates(target addr.Addr, count int, excludes []addr.Addr) (candidates []addr.Addr, err error) {
	ctx := tk.context()
	var excludesStrSlice []string
	for _, x := range excludes {
		excludesStrSlice = append(excludesStrSlice, x.String())
	}
	req := &tracker_interface.GetCandidatesRequest{
		Target:   target.String(),
		Count:    count,
		Excludes: excludesStrSlice,
	}
	var res tracker_interface.GetCandidatesResponse
	err = tk.transport.Call(ctx, "getCandidates", req, &res, well_known_option.WithScope("sys"))
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
