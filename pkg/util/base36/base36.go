package base36

import (
	"strings"

	"github.com/martinlindhe/base36"
)

func Encode(b []byte) string {
	return strings.ToLower(base36.EncodeBytes(b))
}

func Decode(b string) []byte {
	return base36.DecodeToBytes(strings.ToUpper(b))
}
