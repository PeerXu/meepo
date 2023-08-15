package meepo_event_listener

import (
	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

type EventListener interface {
	meepo_eventloop_interface.Handler

	Listen(string, meepo_eventloop_interface.HandleFunc) string
	Unlisten(string)

	Close() error
}

func NewEventListener(opts ...NewEventListenerOption) (EventListener, error) {
	o := option.ApplyWithDefault(DefaultNewEventListenerOption(), opts...)

	queueSize, err := GetQueueSize(o)
	if err != nil {
		return nil, err
	}

	el := &eventListener{
		tree:  NewTree(),
		set:   NewSet(),
		queue: make(chan meepo_eventloop_interface.Event, queueSize),
	}
	go el.start()

	return el, nil
}

type eventListener struct {
	tree  Tree
	set   Set
	queue chan meepo_eventloop_interface.Event
}

func (el *eventListener) start() {
	for evt := range el.queue {
		el.Handle(evt)
	}
}
