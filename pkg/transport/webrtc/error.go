package webrtc_transport

import "fmt"

var (
	GatherTimeoutError = fmt.Errorf("gather timeout")
)

func UnsupportedRoleError(name string) error {
	return fmt.Errorf("Unsupported role %s", name)
}
