package meepo_eventloop_interface

import "time"

type Event interface {
	Name() string
	ID() string
	HappenedAt() time.Time
	Data() Data
	Get(string) any
}
