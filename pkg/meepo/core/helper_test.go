package meepo_core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func TestViewToMap(t *testing.T) {
	tv := sdk_interface.TransportView{
		Addr:    "x",
		Session: "y",
		State:   "z",
	}
	v := viewToMap(tv)
	assert.Equal(t, "x", v["addr"])
	assert.Equal(t, "y", v["session"])
	assert.Equal(t, "z", v["state"])
}
