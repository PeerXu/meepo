package tracker_transport

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type TransportTracker struct {
	transport meepo_interface.Transport
}

func NewTransportTracker(opts ...tracker_core.NewTrackerOption) (tracker_core.Tracker, error) {
	o := option.Apply(opts...)

	transport, err := transport_core.GetTransport(o)
	if err != nil {
		return nil, err
	}

	return &TransportTracker{
		transport: transport,
	}, nil
}

func init() {
	tracker_core.RegisterNewTrackerFunc("transport", NewTransportTracker)
}
