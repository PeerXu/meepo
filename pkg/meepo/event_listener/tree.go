package meepo_event_listener

import (
	"fmt"

	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

type Tree map[string]any

func NewTree() Tree {
	return Tree(make(map[string]any))
}

func (t Tree) SubTree(key string) Tree {
	return t.mustTree(key)
}

func (t Tree) RangeHandleFunc(key string, cb func(id string, fn meepo_eventloop_interface.HandleFunc) bool) {
	t.mustNode(key).Range(cb)
}

func (t Tree) RegHandleFunc(key string, id string, fn meepo_eventloop_interface.HandleFunc) {
	t.mustNode(key).Store(id, fn)
}

func (t Tree) UnregHandleFunc(key string, id string) {
	t.mustNode(key).Delete(id)
}

func (t Tree) joinNodeKey(key string) string {
	return fmt.Sprintf("%s#node", key)
}

func (t Tree) joinTreeKey(key string) string {
	return fmt.Sprintf("%s#tree", key)
}

func (t Tree) mustNode(key string) msync.GenericsMap[string, meepo_eventloop_interface.HandleFunc] {
	nk := t.joinNodeKey(key)
	v, ok := t[nk]
	if !ok {
		v = msync.NewMap[string, meepo_eventloop_interface.HandleFunc]()
		t[nk] = v
	}
	return v.(msync.GenericsMap[string, meepo_eventloop_interface.HandleFunc])
}

func (t Tree) mustTree(key string) Tree {
	tk := t.joinTreeKey(key)
	v, ok := t[tk]
	if !ok {
		v = map[string]any{}
		t[tk] = v
	}
	return v.(map[string]any)
}
