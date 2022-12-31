package crypto_core

import (
	"crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	pubk, prik, _ := ed25519.GenerateKey(nil)
	signer := NewSigner(pubk, prik)
	for _, c := range []struct {
		p *Packet
	}{
		{&Packet{
			Source:      pubk,
			Destination: pubk,
			Nonce:       []byte{0x00},
			CipherText:  []byte{0x00},
		}},
	} {
		assert.Nil(t, signer.Sign(c.p))
		assert.Nil(t, signer.Verify(c.p))
	}
}
