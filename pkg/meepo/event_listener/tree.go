package meepo_event_listener

import (
	"fmt"
	"sync"

	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

// type Tree map[string]any

type Tree struct {
	m *sync.Map
}

func NewTree() *Tree {
	return &Tree{m: &sync.Map{}}
}

func (t *Tree) SubTree(key string) *Tree {
	return t.mustTree(key)
}

func (t *Tree) RangeHandleFunc(key string, cb func(id string, fn meepo_eventloop_interface.HandleFunc) bool) {
	t.mustNode(key).Range(cb)
}

func (t *Tree) RegHandleFunc(key string, id string, fn meepo_eventloop_interface.HandleFunc) {
	t.mustNode(key).Store(id, fn)
}

func (t *Tree) UnregHandleFunc(key string, id string) {
	t.mustNode(key).Delete(id)
}

func (t *Tree) joinNodeKey(key string) string {
	return fmt.Sprintf("%s#node", key)
}

func (t *Tree) joinTreeKey(key string) string {
	return fmt.Sprintf("%s#tree", key)
}

func (t *Tree) mustNode(key string) msync.GenericMap[string, meepo_eventloop_interface.HandleFunc] {
	nk := t.joinNodeKey(key)
	v, ok := t.m.Load(nk)
	if !ok {
		v = msync.NewMap[string, meepo_eventloop_interface.HandleFunc]()
		t.m.Store(nk, v)
	}
	return v.(msync.GenericMap[string, meepo_eventloop_interface.HandleFunc])
}

func (t *Tree) mustTree(key string) *Tree {
	tk := t.joinTreeKey(key)
	v, ok := t.m.Load(tk)
	if !ok {
		v = NewTree()
		t.m.Store(tk, v)
	}
	return v.(*Tree)
}
