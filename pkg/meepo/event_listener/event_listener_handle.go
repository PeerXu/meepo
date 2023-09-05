package meepo_event_listener

import meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"

func (el *eventListener) Handle(e meepo_eventloop_interface.Event) {
	c := DefaultChainParser.Parse(e.Name())
	el.handle(c, el.tree, e)
}

func (el *eventListener) genCallback(key string, t Tree, e meepo_eventloop_interface.Event) func(id string, fn meepo_eventloop_interface.HandleFunc) bool {
	return func(id string, fn meepo_eventloop_interface.HandleFunc) bool {
		if !el.set.Has(id) {
			t.UnregHandleFunc(key, id)
		} else {
			fn(e)
		}
		return true
	}
}

func (el *eventListener) handle(c Chain, t Tree, e meepo_eventloop_interface.Event) {
	head := c.Head()
	rest := c.Rest()

	t.RangeHandleFunc(TAIL_WILDCARD, el.genCallback(TAIL_WILDCARD, t, e))

	if !rest.IsNull() {
		el.handle(rest, t.SubTree(WILDCARD), e)
		el.handle(rest, t.SubTree(head), e)
		return
	}

	t.RangeHandleFunc(WILDCARD, el.genCallback(WILDCARD, t, e))
	t.RangeHandleFunc(head, el.genCallback(head, t, e))
}
