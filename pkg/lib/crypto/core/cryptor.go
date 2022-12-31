package crypto_core

import (
	"crypto/ed25519"
	"io"

	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
)

type Cryptor = crypto_interface.Cryptor

type cryptor struct {
	reader io.Reader
	pubk   ed25519.PublicKey
	prik   ed25519.PrivateKey
}

func NewCryptor(pubk ed25519.PublicKey, prik ed25519.PrivateKey, rd io.Reader) Cryptor {
	return &cryptor{
		reader: rd,
		pubk:   pubk,
		prik:   prik,
	}
}
