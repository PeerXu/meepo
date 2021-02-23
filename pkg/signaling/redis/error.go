package redis_signaling

import "fmt"

func SessionChannelExistError(session int32) error {
	return fmt.Errorf("SessionChannel: %d exist", session)
}

func SessionChannelNotExistError(session int32) error {
	return fmt.Errorf("SessionChannel: %d not exist", session)
}

func SessionChannelClosedError(session int32) error {

	return fmt.Errorf("session channel(%d) closed", session)
}

var (
	WaitWiredEventTimeoutError   = fmt.Errorf("wait wired event timeout")
	NotAvailableRedisClientError = fmt.Errorf("not available redis client")
)
