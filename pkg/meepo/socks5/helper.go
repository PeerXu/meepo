package meepo_socks5

import (
	"strings"

	"github.com/PeerXu/meepo/pkg/lib/addr"
)

func parseDomain(str string, root string, addrSize int) (sinkHost string, target string, err error) {
	strSize := len(str)
	rootSize := len(root)

	if strSize < rootSize+addrSize {
		return "", "", ErrInvalidDomainFn(str)
	}

	if !strings.HasSuffix(str, root) {
		return "", "", ErrInvalidDomainFn(str)
	}

	target = str[strSize-rootSize-addrSize : strSize-rootSize]
	if strSize > rootSize+addrSize {
		sinkHost = str[:strSize-rootSize-addrSize-1]
	} else {
		sinkHost = "127.0.0.1"
	}

	return
}

func (ss *Socks5Server) parseDomain(str string) (sinkHost string, target addr.Addr, err error) {
	sinkHost, targetStr, err := parseDomain(str, ss.root, addr.ADDR_STR_SIZE)
	if err != nil {
		return "", "", err
	}

	target, err = addr.FromString(targetStr)
	if err != nil {
		return "", "", err
	}

	return
}
