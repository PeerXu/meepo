package meepo_interface

import (
	"context"
	"net"
)

type Teleportation interface {
	ID() string
	Addr() Addr
	SourceAddr() net.Addr
	SinkAddr() net.Addr
	Close(context.Context) error
	Mode() string
}

type Teleportationer interface {
	NewTeleportation(ctx context.Context, transportAddr Addr, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress string, opts ...NewTeleportationOption) (Teleportation, error)
	ListTeleportations(ctx context.Context, opts ...ListTeleportationsOption) ([]Teleportation, error)
	GetTeleportation(ctx context.Context, id string, opts ...GetTeleportationOption) (Teleportation, error)
}
