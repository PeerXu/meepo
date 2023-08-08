package meepo_eventloop_interface

type Event interface {
	Name() string
	ID() string
	Data() Data
	Get(string) any
}
