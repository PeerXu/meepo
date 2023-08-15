package meepo_event_listener

import meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"

func (el *eventListener) Handle(e meepo_eventloop_interface.Event) {
	c := DefaultChainParser.Parse(e.Name())
	el.handle(c, el.tree, e)
}

func (el *eventListener) handle(c Chain, t Tree, e meepo_eventloop_interface.Event) {
	head := c.Head()
	rest := c.Rest()

	if !rest.IsNull() {
		el.handle(rest, t.SubTree("*"), e)
		el.handle(rest, t.SubTree(head), e)
		return
	}

	t.RangeHandleFunc("*", func(id string, fn meepo_eventloop_interface.HandleFunc) bool {
		if !el.set.Has(id) {
			t.UnregHandleFunc("*", id)
		} else {
			fn(e)
		}
		return true
	})

	t.RangeHandleFunc(head, func(id string, fn meepo_eventloop_interface.HandleFunc) bool {
		if !el.set.Has(id) {
			t.UnregHandleFunc(head, id)
		} else {
			fn(e)
		}
		return true
	})
}
