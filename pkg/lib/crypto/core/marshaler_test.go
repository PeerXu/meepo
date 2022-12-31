package crypto_core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaler(t *testing.T) {
	for _, c := range []struct {
		p *Packet
	}{
		{&Packet{
			Source:      []byte{0x01},
			Destination: []byte{0x02},
			Nonce:       []byte{0x03},
			CipherText:  []byte{0x04},
			Signature:   []byte{0x05},
		}},
	} {
		buf, err := MarshalPacket(c.p)
		assert.Nil(t, err)
		p, err := UnmarshalPacket(buf)
		assert.Nil(t, err)
		assert.Equal(t, c.p, p)
	}
}
