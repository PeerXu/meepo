package crypto_core

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
)

type (
	EncryptOption = crypto_interface.EncryptOption
	DecryptOption = crypto_interface.DecryptOption
)

const (
	OPTION_SIGNER  = "signer"
	OPTION_CRYPTOR = "cryptor"
)

func WithSigner(x Signer) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_SIGNER] = x
	}
}

func GetSigner(o option.Option) (Signer, error) {
	var x Signer
	i := o.Get(OPTION_SIGNER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_SIGNER)
	}

	v, ok := i.(Signer)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(x, i)
	}

	return v, nil
}

func WithCryptor(x Cryptor) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_CRYPTOR] = x
	}
}

func GetCryptor(o option.Option) (Cryptor, error) {
	var x Cryptor
	i := o.Get(OPTION_CRYPTOR).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_CRYPTOR)
	}

	v, ok := i.(Cryptor)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(x, i)
	}

	return v, nil

}
