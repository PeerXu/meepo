package meepo_eventloop_core

import meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"

func (el *eventLoop) AddHandler(h meepo_eventloop_interface.Handler) string {
	id := randomID(HANDLER_ID_SIZE)
	el.handlers.Store(id, h)
	return id
}
