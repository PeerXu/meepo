package signaling

import "fmt"

var (
	WireTimeoutError        = fmt.Errorf("Wire timeout")
	NextHopUnreachableError = fmt.Errorf("Next hop unreachable")
)

func UnsupportedSignalingEngine(name string) error {
	return fmt.Errorf("Unsupported signaling engine: %s", name)
}
