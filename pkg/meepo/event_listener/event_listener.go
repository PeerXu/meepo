package meepo_event_listener

import meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"

type EventListener interface {
	meepo_eventloop_interface.Handler

	Listen(string, meepo_eventloop_interface.HandleFunc) string
	Unlisten(string)
}

func NewEventListener() EventListener {
	return &eventListener{
		t: NewTree(),
		s: NewSet(),
	}
}

type eventListener struct {
	t Tree
	s Set
}
