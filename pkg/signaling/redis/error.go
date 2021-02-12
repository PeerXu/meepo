package redis_signaling

import "fmt"

func ErrReregisterSessionChannel(session int32) error {
	return fmt.Errorf("reregister session channel: %d", session)
}

func ErrUnregisteredSessionChannel(session int32) error {
	return fmt.Errorf("unregistered session channel: %d", session)
}

func ErrSessionChannelClosed(session int32) error {

	return fmt.Errorf("session channel(%d) closed", session)
}

var (
	ErrWaitWiredEventTimeout   = fmt.Errorf("wait wired event timeout")
	ErrNotAvailableRedisClient = fmt.Errorf("not available redis client")
)
