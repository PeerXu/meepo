package dialer_interface

import "context"

type Dialer interface {
	Dial(ctx context.Context, network string, address string, opts ...DialOption) (Conn, error)
}
