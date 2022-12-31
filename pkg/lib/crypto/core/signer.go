package crypto_core

import (
	"crypto/ed25519"

	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
)

type Signer = crypto_interface.Signer

type signer struct {
	pubk ed25519.PublicKey
	prik ed25519.PrivateKey
}

func NewSigner(pubk ed25519.PublicKey, prik ed25519.PrivateKey) Signer {
	return &signer{
		pubk: pubk,
		prik: prik,
	}
}
