package crypt

import (
	"crypto/ed25519"

	"github.com/teserakt-io/golang-ed25519/extra25519"
	"golang.org/x/crypto/curve25519"
)

func ed25519PublicKeyToCurve25519(pubk ed25519.PublicKey) [32]byte {
	var curve [32]byte
	var pubk32 [32]byte
	copy(pubk32[:], pubk[:32])
	extra25519.PublicKeyToCurve25519(&curve, &pubk32)
	return curve
}

func ed25519PrivateKeyToCurve25519(prik ed25519.PrivateKey) [32]byte {
	var curve [32]byte
	var prik64 [64]byte
	copy(prik64[:], prik[:64])
	extra25519.PrivateKeyToCurve25519(&curve, &prik64)
	return curve
}

func CalcSharedSecret(pubk ed25519.PublicKey, prik ed25519.PrivateKey) [32]byte {
	var secret [32]byte
	curvePubk := ed25519PublicKeyToCurve25519(pubk)
	curvePrik := ed25519PrivateKeyToCurve25519(prik)
	curve25519.ScalarMult(&secret, &curvePrik, &curvePubk)
	return secret
}
