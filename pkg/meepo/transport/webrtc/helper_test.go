package transport_webrtc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResponseSession(t *testing.T) {
	for i, c := range []struct {
		s      string
		expect string
	}{
		{"ffffffff", "fffffffe"},
		{"f3c594d3", "f3c594d2"},
	} {
		assert.Equal(t, c.expect, parseResponseSession(c.s), "i=%v, s=%v", i, c.s)
	}
}

func TestFailed(t *testing.T) {
	go func() {
		select {}
	}()
}
