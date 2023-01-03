package tracker_core

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrUnsupportedTracker, ErrUnsupportedTrackerFn = errors.NewErrorAndErrorFunc[string]("unsupported tracker")
)
