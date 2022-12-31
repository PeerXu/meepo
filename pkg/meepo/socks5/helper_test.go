package meepo_socks5

import (
	"testing"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/stretchr/testify/assert"
)

func TestParseDomain(t *testing.T) {
	for _, c := range []struct {
		str            string
		root           string
		addrSize       int
		expectSinkHost string
		expectTarget   string
	}{
		{"610dbbx70vrbqotu3w2pij6l8wp1bif9r3fhb5b5dbene07bgka.mpo", ".mpo", addr.ADDR_STR_SIZE, "127.0.0.1", "610dbbx70vrbqotu3w2pij6l8wp1bif9r3fhb5b5dbene07bgka"},
		{"a.b", ".b", 1, "127.0.0.1", "a"},
		{"192.168.1.199.a.b", ".b", 1, "192.168.1.199", "a"},
		{".a.b", ".b", 1, "", "a"},
	} {
		sinkHost, target, err := parseDomain(c.str, c.root, c.addrSize)
		assert.Nil(t, err)
		assert.Equal(t, c.expectSinkHost, sinkHost)
		assert.Equal(t, c.expectTarget, target)
	}
}
