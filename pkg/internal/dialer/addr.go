package dialer

import "net"

type addr struct {
	network string
	address string
}

func (x *addr) Network() string {
	return x.network
}

func (x *addr) String() string {
	return x.address
}

func NewAddr(network, address string) net.Addr {
	return &addr{network: network, address: address}
}
