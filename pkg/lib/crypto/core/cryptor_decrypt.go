package crypto_core

func (x *cryptor) Decrypt(p *Packet, opts ...DecryptOption) ([]byte, error) {
	secret := SecretFromEd25519(p.Nonce, x.prik)
	gcm, err := NewGCM(secret)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, p.Nonce[:12], p.CipherText, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
