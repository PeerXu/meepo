package tracker_core

import "github.com/PeerXu/meepo/pkg/internal/errors"

var (
	ErrUnsupportedTracker, ErrUnsupportedTrackerFn = errors.NewErrorAndErrorFunc[string]("unsupported tracker")
)
