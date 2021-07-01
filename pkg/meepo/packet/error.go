package packet

import "fmt"

var (
	ErrInvalidPacket = fmt.Errorf("invalid packet")
	ErrPacketIsNil   = fmt.Errorf("packet is nil")
)
