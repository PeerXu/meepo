package crypto_interface

type Signer interface {
	Sign(*Packet) error
	Verify(*Packet) error
}
