package crypto_core

import (
	"bytes"
)

func (x *marshaler) Marshal(p *Packet) (out []byte, err error) {
	bb := new(bytes.Buffer)

	if _, err = bb.Write(x.magicCode); err != nil {
		return
	}

	for _, buf := range [][]byte{p.Source, p.Destination, p.Nonce, p.CipherText, p.Signature} {
		if err = WriteBufferWithSize(bb, buf); err != nil {
			return
		}
	}

	return bb.Bytes(), nil
}
