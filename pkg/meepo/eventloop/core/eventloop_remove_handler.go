package meepo_eventloop_core

func (el *eventLoop) RemoveHandler(id string) {
	el.handlers.Delete(id)
}
