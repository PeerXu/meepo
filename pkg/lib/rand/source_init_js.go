//go:build js && wasm

package rand

import (
	"math/rand"
	"time"
)

func init() {
	globalSource = rand.NewSource(time.Now().UnixNano())
}
