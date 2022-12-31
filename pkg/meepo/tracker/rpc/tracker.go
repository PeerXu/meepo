package tracker_rpc

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

type RPCTracker struct {
	addr   addr.Addr
	caller rpc_core.Caller
}

func NewRPCTracker(opts ...tracker_core.NewTrackerOption) (tracker_core.Tracker, error) {
	o := option.Apply(opts...)

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	caller, err := rpc_core.GetCaller(o)
	if err != nil {
		return nil, err
	}

	return &RPCTracker{
		addr:   addr,
		caller: caller,
	}, nil
}

func init() {
	tracker_core.RegisterNewTrackerFunc("rpc", NewRPCTracker)
}
