package transport

import (
	"io"
)

type DataChannelState int

const (
	DataChannelStateConnecting DataChannelState = iota + 1
	DataChannelOpen
	DataChannelClosing
	DataChannelClosed
)

var (
	DataChannelStateStr = []string{
		DataChannelStateConnecting: "connecting",
		DataChannelOpen:            "open",
		DataChannelClosing:         "closing",
		DataChannelClosed:          "closed",
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
