package crypto_core

import (
	"crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptor(t *testing.T) {
	pubk, prik, _ := ed25519.GenerateKey(nil)
	cryptor := NewCryptor(pubk, prik, nil)
	for _, c := range []struct {
		plaintext []byte
	}{
		{[]byte("hello, world")},
	} {
		p, err := cryptor.Encrypt(pubk, c.plaintext)
		assert.Nil(t, err)
		plaintext, err := cryptor.Decrypt(p)
		assert.Nil(t, err)
		assert.Equal(t, c.plaintext, plaintext)
	}
}
