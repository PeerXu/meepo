package tracker_core

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	OPTION_TRACKERS = "trackers"

	METHOD_GET_CANDIDATES      = "getCandidates"
	METHOD_NEW_TRANSPORT       = "newTransport"
	METHOD_ADD_PEER_CONNECTION = "addPeerConnection"
)

type NewTrackerOption = option.ApplyOption

func WithTrackers(ts ...Tracker) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_TRACKERS] = ts
	}
}

func GetTrackers(o option.Option) ([]Tracker, error) {
	i := o.Get(OPTION_TRACKERS).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_TRACKERS)
	}
	v, ok := i.([]Tracker)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}
	return v, nil
}
