package meepo_eventloop_core

import (
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
)

const (
	BASIC_EVENT_ID_SIZE = 8
)

type basicEvent struct {
	name string
	id   string
	data meepo_eventloop_interface.Data
}

func (e *basicEvent) Name() string {
	return e.name
}

func (e *basicEvent) ID() string {
	return e.id
}

func (e *basicEvent) Data() meepo_eventloop_interface.Data {
	return e.data
}

func (e *basicEvent) Get(key string) any {
	if e.data == nil {
		return nil
	}
	return e.data[key]
}

func NewEvent(name string, data meepo_eventloop_interface.Data) meepo_eventloop_interface.Event {
	return &basicEvent{
		name: name,
		id:   randomID(BASIC_EVENT_ID_SIZE),
		data: data,
	}
}
