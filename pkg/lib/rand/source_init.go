//go:build !(js && wasm)

package rand

import "math/rand"

func init() {
	globalSource = rand.NewSource(rand.Int63())
}
