package dialer

import "net"

type addr = net.UnixAddr

func NewAddr(network, address string) net.Addr {
	return &addr{Net: network, Name: address}
}
