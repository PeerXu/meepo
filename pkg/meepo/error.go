package meepo

import (
	"fmt"
	"net"
)

var (
	GatherTimeoutError = fmt.Errorf("Gather timeout")
)

func NotListenableAddressError(addr net.Addr) error {
	return fmt.Errorf("Not listenable address: %s %s", addr.Network(), addr.String())
}
