package meepo_interface

import "context"

type Transporter interface {
	NewTransport(context.Context, Addr, ...NewTransportOption) (Transport, error)
	ListTransports(context.Context, ...ListTransportsOption) ([]Transport, error)
	GetTransport(context.Context, Addr, ...GetTransportOption) (Transport, error)
}

type Transport interface {
	Addr() Addr
	Close(context.Context) error
	WaitReady() error
	State() TransportState

	Channeler
	Caller
	Handler
}
