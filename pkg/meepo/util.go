package meepo

import (
	"net"
)

func checkAddrIsListenable(addr net.Addr) (net.Addr, error) {
	lis, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		return nil, err
	}
	defer lis.Close()

	tcpAddr, _ := net.ResolveTCPAddr(lis.Addr().Network(), lis.Addr().String())

	return tcpAddr, nil
}

// TODO(Peer): More robustness
func getListenableAddr() net.Addr {
	for {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
		if addr, err := checkAddrIsListenable(addr); err == nil {
			return addr
		}
	}
}
