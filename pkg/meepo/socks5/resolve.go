package meepo_socks5

import (
	"context"
	"net"
)

func (ss *Socks5Server) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	return ctx, net.IPv4(127, 0, 0, 1), nil
}
