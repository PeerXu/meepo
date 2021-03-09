package sync

import "fmt"

var (
	ChannelExistError    = fmt.Errorf("Channel exist")
	ChannelNotExistError = fmt.Errorf("Channel not exist")
)
