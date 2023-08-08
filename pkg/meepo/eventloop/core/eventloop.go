package meepo_eventloop_core

import (
	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

const (
	HANDLER_ID_SIZE = 8
)

type eventLoop struct {
	handlers msync.GenericsMap[string, meepo_eventloop_interface.Handler]
}

func NewEventLoop() meepo_eventloop_interface.EventLoop {
	return &eventLoop{
		handlers: msync.NewMap[string, meepo_eventloop_interface.Handler](),
	}
}
