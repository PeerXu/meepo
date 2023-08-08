package meepo_eventloop_core

import meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"

func (el *eventLoop) Emit(e meepo_eventloop_interface.Event) {
	el.handlers.Range(func(_ string, h meepo_eventloop_interface.Handler) bool {
		h.Handle(e)
		return true
	})
}
