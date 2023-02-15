package meepo_core

import crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"

func (mp *Meepo) encryptMessage(target Addr, in any) (*crypto_core.Packet, error) {
	buf, err := mp.marshaler.Marshal(in)
	if err != nil {
		return nil, err
	}

	out, err := mp.cryptor.Encrypt(target.Bytes(), buf)
	if err != nil {
		return nil, err
	}

	if err = mp.signer.Sign(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (mp *Meepo) decryptMessage(in *crypto_core.Packet, out any) error {
	buf, err := mp.cryptor.Decrypt(in)
	if err != nil {
		return err
	}

	if err = mp.unmarshaler.Unmarshal(buf, out); err != nil {
		return err
	}

	return nil
}
