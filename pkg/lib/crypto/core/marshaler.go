package crypto_core

var (
	MAGIC_CODE = []byte{0x3E, 0xE0}
)

type marshaler struct {
	magicCode []byte
}

var DefaultMarshaler = &marshaler{magicCode: MAGIC_CODE[:]}

func MarshalPacket(p *Packet) ([]byte, error) {
	return DefaultMarshaler.Marshal(p)
}

func UnmarshalPacket(b []byte) (*Packet, error) {
	return DefaultMarshaler.Unmarshal(b)
}
