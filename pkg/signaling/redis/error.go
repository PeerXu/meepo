package redis_signaling

import "fmt"

var (
	NotAvailableRedisClientError = fmt.Errorf("Not available redis client")
	SessionChannelClosedError    = fmt.Errorf("Session channel closed")
	SessionChannelNotExistError  = fmt.Errorf("Session channel not exist")
	SessionChannelExistError     = fmt.Errorf("Session channel exist")
)
