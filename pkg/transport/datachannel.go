package transport

import (
	"io"
)

type DataChannelState int

const (
	DataChannelStateConnecting DataChannelState = iota + 1
	DataChannelStateOpen
	DataChannelStateClosing
	DataChannelStateClosed
)

var (
	DataChannelStateStr = []string{
		DataChannelStateConnecting: "connecting",
		DataChannelStateOpen:       "open",
		DataChannelStateClosing:    "closing",
		DataChannelStateClosed:     "closed",
	}
)

func (t DataChannelState) String() string {
	return DataChannelStateStr[t]
}

type DataChannel interface {
	Transport() Transport
	Label() string
	State() DataChannelState
	OnOpen(func())
	io.ReadWriteCloser
}
