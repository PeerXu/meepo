package crypto_interface

type Cryptor interface {
	Encrypt(dest, text []byte, opts ...EncryptOption) (*Packet, error)
	Decrypt(*Packet, ...DecryptOption) ([]byte, error)
}
