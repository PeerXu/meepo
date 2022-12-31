package cpl_test

import (
	"testing"

	"github.com/PeerXu/meepo/pkg/internal/cpl"
	"github.com/stretchr/testify/assert"
)

func TestIndexOfLongestCommonPrefix(t *testing.T) {
	for i, c := range []struct {
		x, y []byte
		lcp  int
	}{
		{[]byte{0x00}, []byte{0x00}, 8},
		{[]byte{0x00}, []byte{0x01}, 7},
		{[]byte{0x00}, []byte{0x02}, 6},
		{[]byte{0x00}, []byte{0x04}, 5},
		{[]byte{0x00}, []byte{0x08}, 4},
		{[]byte{0x00}, []byte{0x10}, 3},
		{[]byte{0x00}, []byte{0x20}, 2},
		{[]byte{0x00}, []byte{0x40}, 1},
		{[]byte{0x80}, []byte{0x00}, 0},
	} {
		cx, cy := cpl.FromBytes(c.x), cpl.FromBytes(c.y)
		idxOfLcp := cpl.CommonPrefixLen(cx, cy)
		assert.Equal(t, c.lcp, idxOfLcp, "index=%v x=%v y=%v", i, c.x, c.y)
	}
}
