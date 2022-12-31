package crypto_core

import (
	"crypto/ed25519"
)

func (x *cryptor) Encrypt(dest, data []byte, opts ...EncryptOption) (*Packet, error) {
	randPubk, randPrik, err := ed25519.GenerateKey(x.reader)
	if err != nil {
		return nil, err
	}

	destPubk := ed25519.PublicKey(dest)
	nonce := randPubk[:12]
	secret := SecretFromEd25519(destPubk, randPrik)
	gcm, err := NewGCM(secret)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return &Packet{
		Source:      []byte(x.pubk),
		Destination: dest,
		Nonce:       randPubk,
		CipherText:  ciphertext,
	}, nil
}
