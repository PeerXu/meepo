package addr

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/PeerXu/meepo/pkg/lib/base36"
)

const (
	MAGIC_CODE    byte = 0x22
	ADDR_SIZE          = ed25519.PublicKeySize + 1
	ADDR_STR_SIZE      = 51
)

type Addr string

func (x Addr) String() string {
	return string(x)
}

func (x Addr) Bytes() []byte {
	return base36.Decode(x.String())[1:]
}

func (x Addr) Equal(y Addr) bool {
	return x.String() == y.String()
}

func FromString(x string) (Addr, error) {
	return FromBytes(base36.Decode(x))
}

func FromBytes(x []byte) (Addr, error) {
	if len(x) != ADDR_SIZE {
		return "", ErrInvalidAddrStringFn(base64.RawStdEncoding.EncodeToString(x))
	}

	if x[0] != MAGIC_CODE {
		return "", ErrInvalidAddrStringFn(base64.RawStdEncoding.EncodeToString(x))
	}
	return Addr(base36.Encode(x)), nil
}

func FromBytesWithoutMagicCode(x []byte) (Addr, error) {
	if len(x) != ADDR_SIZE-1 {
		return "", ErrInvalidAddrStringFn(base64.RawStdEncoding.EncodeToString(x))
	}

	return FromBytes(append([]byte{MAGIC_CODE}, x...))
}

func Must(x Addr, _ error) Addr {
	return x
}
