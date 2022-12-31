package meepo_interface

import (
	"context"
	"net"
)

type Channel interface {
	ID() uint16
	Mode() string
	WaitReady() error
	State() ChannelState
	Conn() Conn
	Close(context.Context) error
	IsSource() bool
	IsSink() bool
	SinkAddr() net.Addr
}

type Channeler interface {
	NewChannel(ctx context.Context, network string, address string, opts ...NewChannelOption) (Channel, error)
	ListChannels(ctx context.Context, opts ...ListChannelsOption) ([]Channel, error)
	GetChannel(ctx context.Context, id uint16) (Channel, error)
}
