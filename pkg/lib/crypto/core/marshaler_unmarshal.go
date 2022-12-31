package crypto_core

import (
	"bytes"
)

func (x *marshaler) Unmarshal(b []byte) (p *Packet, err error) {
	var pp Packet

	mcSize := len(x.magicCode)
	if !bytes.Equal(x.magicCode, b[:mcSize]) {
		return nil, ErrInvalidBuffer
	}

	rd := bytes.NewReader(b[mcSize:])
	for _, ptr := range []*[]byte{&pp.Source, &pp.Destination, &pp.Nonce, &pp.CipherText, &pp.Signature} {
		if err = ReadBufferWithSize(rd, ptr); err != nil {
			return
		}
	}

	return &pp, nil
}
