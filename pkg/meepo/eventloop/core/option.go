package meepo_eventloop_core

import (
	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

const (
	OPTION_EVENTLOOP = "eventLoop"
)

var (
	WithEventLoop, GetEventLoop = option.New[meepo_eventloop_interface.EventLoop](OPTION_EVENTLOOP)
)
