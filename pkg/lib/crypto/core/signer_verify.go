package crypto_core

import (
	"crypto/ed25519"
)

func (x *signer) Verify(p *Packet) error {
	q := *p
	s := p.Signature
	q.Signature = make([]byte, ed25519.SignatureSize)
	msg, err := MarshalPacket(&q)
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(q.Source), msg, s) {
		return ErrInvalidSignature
	}
	return nil
}
