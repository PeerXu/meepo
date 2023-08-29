package meepo_interface

type ChannelState string

const (
	CHANNEL_STATE_NEW        ChannelState = "new"
	CHANNEL_STATE_CONNECTING ChannelState = "connecting"
	CHANNEL_STATE_OPEN       ChannelState = "open"
	CHANNEL_STATE_CLOSING    ChannelState = "closing"
	CHANNEL_STATE_CLOSED     ChannelState = "closed"
)

func (x ChannelState) String() string {
	return string(x)
}
