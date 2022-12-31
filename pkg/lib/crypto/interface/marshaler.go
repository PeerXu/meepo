package crypto_interface

type Marshaler interface {
	Marshal(*Packet) ([]byte, error)
	Unmarshal([]byte) (*Packet, error)
}
