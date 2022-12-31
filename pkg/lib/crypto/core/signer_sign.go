package crypto_core

import (
	"crypto/ed25519"
)

func (x *signer) Sign(p *Packet) error {
	p.Signature = make([]byte, ed25519.SignatureSize)
	msg, err := MarshalPacket(p)
	if err != nil {
		return err
	}
	p.Signature = ed25519.Sign(x.prik, msg)
	return nil
}
