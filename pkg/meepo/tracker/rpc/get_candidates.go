package tracker_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	lib_protocol "github.com/PeerXu/meepo/pkg/lib/protocol"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

func (tk *RPCTracker) GetCandidates(target addr.Addr, requests int, excludes []addr.Addr) (candidates []addr.Addr, err error) {
	ctx := tk.context()
	var excludesStrSlice []string
	for _, x := range excludes {
		excludesStrSlice = append(excludesStrSlice, x.String())
	}

	req := &tracker_interface.GetCandidatesRequest{
		Protocol: lib_protocol.VERSION.String(),
		Target:   target.String(),
		Requests: requests,
		Excludes: excludesStrSlice,
	}
	var res tracker_interface.GetCandidatesResponse
	err = tk.caller.Call(ctx, tracker_core.METHOD_GET_CANDIDATES, req, &res, well_known_option.WithDestination(tk.addr.Bytes()))
	if err != nil {
		return
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
