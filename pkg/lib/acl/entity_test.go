package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, c := range []struct {
		s             string
		expectID      string
		expectNetwork string
		expectHost    string
		expectPort    string
	}{
		{"a,tcp,192.168.1.1:8080", "a", "tcp", "192.168.1.1", "8080"},
		{"a,tcp,*:8080", "a", "tcp", "*", "8080"},
		{"*,*,127.0.0.1:8080", "*", "*", "127.0.0.1", "8080"},
		{"a,tcp,*", "a", "tcp", "*", "*"},
		{"a,*,*", "a", "*", "*", "*"},
		{"a,*", "a", "*", "*", "*"},
		{"a,tcp", "a", "tcp", "*", "*"},
		{"a", "a", "*", "*", "*"},
		{"*", "*", "*", "*", "*"},
	} {
		e, err := Parse(c.s)
		assert.Nil(t, err)
		assert.Equal(t, c.expectID, e.ID())
		assert.Equal(t, c.expectNetwork, e.Network())
		assert.Equal(t, c.expectHost, e.Host())
		assert.Equal(t, c.expectPort, e.Port())
	}
}

func TestContains(t *testing.T) {
	for _, c := range []struct {
		s          string
		challenges []string
		expect     bool
	}{
		{"*", []string{"a,tcp,192.168.199.1:8080"}, true},
		{"a", []string{"a,tcp,10.1.1.1:22"}, true},
		{"*,tcp,192.168.100.1:8080", []string{"a,tcp,192.168.199.1:8080"}, false},
		{"*,tcp,10.1.1.0/24:8080", []string{"a,tcp,10.1.1.233:8080"}, true},
		{"*,tcp,10.1.1.0/24:8080", []string{
			"a,tcp,10.2.1.233:8080",
			"a,tcp,www.baidu.com:8080",
		}, false},
	} {
		e, err := Parse(c.s)
		require.Nil(t, err)

		for _, challengeStr := range c.challenges {
			challenge, err := Parse(challengeStr)
			require.Nil(t, err)
			assert.Equal(t, c.expect, e.Contains(challenge))
		}
	}
}
