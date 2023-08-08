package meepo_eventloop_interface

type HandleFunc func(Event)

func (fn HandleFunc) Handle(e Event) {
	fn(e)
}

type Handler interface {
	Handle(Event)
}
