package meepo_event_listener

import (
	mrand "github.com/PeerXu/meepo/pkg/lib/rand"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

func (el *eventListener) Listen(name string, fn meepo_eventloop_interface.HandleFunc) string {
	c := DefaultChainParser.Parse(name)
	return listen(c, el.tree, el.set, fn)
}

func listen(c Chain, t *Tree, s Set, fn meepo_eventloop_interface.HandleFunc) string {
	var key string

	chainHead := c.Head()
	chainRest := c.Rest()

	if !chainRest.IsNull() {
		return listen(chainRest, t.SubTree(chainHead), s, fn)
	}

	if chainHead == WILDCARD {
		key = TAIL_WILDCARD
	} else {
		key = chainHead
	}

	id := mrand.DefaultStringGenerator.Generate(8)
	t.RegHandleFunc(key, id, fn)
	s.Add(id)

	return id
}
