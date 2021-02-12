package signaling

import "fmt"

func UnsupportedSignalingEngine(name string) error {
	return fmt.Errorf("unsupported signaling engine: %s", name)
}
