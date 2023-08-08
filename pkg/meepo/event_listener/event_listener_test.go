package meepo_event_listener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	meepo_event_listener "github.com/PeerXu/meepo/pkg/meepo/event_listener"
	meepo_eventloop_core "github.com/PeerXu/meepo/pkg/meepo/eventloop/core"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

func TestEventListener(t *testing.T) {
	cnt := 0
	el := meepo_event_listener.NewEventListener()
	f := func(e meepo_eventloop_interface.Event) {
		assert.Equal(t, e.Name(), "x.y.z")
		cnt++
	}

	id1 := el.Listen("x.*.z", f)
	id2 := el.Listen("*.y.z", f)
	el.Handle(meepo_eventloop_core.NewEvent("x.y.z", nil))
	el.Unlisten(id1)
	el.Unlisten(id2)
	el.Handle(meepo_eventloop_core.NewEvent("x.y.z", nil))
	assert.Equal(t, 2, cnt)
}
