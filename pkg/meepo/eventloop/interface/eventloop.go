package meepo_eventloop_interface

type EventLoop interface {
	AddHandler(Handler) string
	RemoveHandler(string)
	Emit(Event)
}
