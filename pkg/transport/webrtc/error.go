package webrtc_transport

import "fmt"

var (
	GatherTimeoutError                = fmt.Errorf("gather timeout")
	WaitDataChannelOpenedTimeoutError = fmt.Errorf("wait data channel opened timeout")
)

func UnsupportedRoleError(name string) error {
	return fmt.Errorf("Unsupported role %s", name)
}
