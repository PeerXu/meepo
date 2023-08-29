package meepo_event_listener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	meepo_event_listener "github.com/PeerXu/meepo/pkg/meepo/event_listener"
	meepo_eventloop_core "github.com/PeerXu/meepo/pkg/meepo/eventloop/core"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

func TestEventListener0(t *testing.T) {
	cnt := 0
	el, err := meepo_event_listener.NewEventListener()
	require.Nil(t, err)

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

func TestEventListener(t *testing.T) {
	for _, st := range []struct {
		rules  []string
		events []struct {
			name string
			skip bool
		}
	}{
		// TEMPLATE
		{[]string{}, []struct {
			name string
			skip bool
		}{}},

		{[]string{
			"x.*.z",
		}, []struct {
			name string
			skip bool
		}{
			{"x.y.z", false},
			{"x.a.z", false},
			{"x.y.z.0", true},
			{"a.y.z", true},
		}},

		{[]string{
			"x.*",
		}, []struct {
			name string
			skip bool
		}{
			{"x.y", false},
			{"x.z", false},
			{"x.y.z", true},
			{"a.b.c", true},
		}},

		{[]string{
			"x",
		}, []struct {
			name string
			skip bool
		}{
			{"x", false},
			{"y", true},
		},
		},
		{[]string{
			"a",
			"b",
			"c",
		}, []struct {
			name string
			skip bool
		}{
			{"a", false},
			{"b", false},
			{"c", false},
			{"d", true},
		}},
		{[]string{
			"x.*",
		}, []struct {
			name string
			skip bool
		}{
			{"x.y", false},
			{"x.z", false},
			{"a", true},
			{"a.b", true},
		}},
	} {
		lis, err := meepo_event_listener.NewEventListener()
		require.Nil(t, err)

		var expects []int
		for i := 0; i < len(st.events); i++ {
			if !st.events[i].skip {
				expects = append(expects, i)
			}
		}

		var actuals []int
		for _, r := range st.rules {
			lis.Listen(r, func(evt meepo_eventloop_interface.Event) {
				actuals = append(actuals, evt.Get("index").(int))
			})
		}

		for i := 0; i < len(st.events); i++ {
			lis.Handle(meepo_eventloop_core.NewEvent(st.events[i].name, map[string]any{"index": i}))
		}

		assert.Equal(t, expects, actuals)
	}
}
