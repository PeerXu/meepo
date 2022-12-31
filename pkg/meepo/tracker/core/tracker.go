package tracker_core

import (
	"sync"

	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

type Tracker = tracker_interface.Tracker

type NewTrackerFunc func(...NewTrackerOption) (Tracker, error)

var newTrackerFuncs sync.Map

func NewTracker(name string, opts ...NewTrackerOption) (Tracker, error) {
	v, ok := newTrackerFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedTrackerFn(name)
	}
	return v.(NewTrackerFunc)(opts...)
}

func RegisterNewTrackerFunc(name string, fn NewTrackerFunc) {
	newTrackerFuncs.Store(name, fn)
}
