package meepo_interface

import "context"

type Dialer interface {
	Dial(ctx context.Context, target Addr, network string, address string, opts ...DialOption) (Conn, error)
}
