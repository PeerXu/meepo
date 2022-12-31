package crypto_core

import (
	"crypto/ed25519"
)

func (x *signer) Verify(p *Packet) error {
	sig := p.Signature
	defer func() {
		p.Signature = sig
	}()
	p.Signature = make([]byte, ed25519.SignatureSize)
	msg, err := MarshalPacket(p)
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(p.Source), msg, sig) {
		return ErrInvalidSignature
	}
	return nil
}
