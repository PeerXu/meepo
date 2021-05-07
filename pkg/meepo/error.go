package meepo

import (
	"fmt"
)

var (
	UnsupportedRequestHandlerError = fmt.Errorf("Unsupported request method")
	TransportNotExistError         = fmt.Errorf("Transport not exists")
	TransportExistError            = fmt.Errorf("Transport already exists")
	TeleportationNotExistError     = fmt.Errorf("Teleportation not exists")
	UnexpectedMessageError         = fmt.Errorf("Unexpected message")
	WaitResponseTimeoutError       = fmt.Errorf("Wait response timeout")
	NotWirableError                = fmt.Errorf("Not wirable")
	HopIsZeroError                 = fmt.Errorf("Hop is zero")
	ReachTransportEdgeError        = fmt.Errorf("Reach transport edge")
	UnsupportedSocks5CommandError  = fmt.Errorf("Unsupported socks5 command")
	NetworkUnreachableError        = fmt.Errorf("Network unreachable")
	UnsupportedNetworkTypeError    = fmt.Errorf("Unsupported network type")
)

func SessionChannelExistError(session int32) error {
	return fmt.Errorf("SessionChannel: %d exist", session)
}

func SessionChannelNotExistError(session int32) error {
	return fmt.Errorf("SessionChannel: %d not exist", session)
}

func SessionChannelClosedError(session int32) error {
	return fmt.Errorf("SessionChannel: %d closed", session)
}

func UnsupportedMessageDecodeDriverError(messageIdentifier string) error {
	return fmt.Errorf("Unsupported message decode driver: %s", messageIdentifier)
}

type sendMessageError struct {
	err error
}

func (t sendMessageError) Error() string {
	return t.err.Error()
}

func SendMessageError(err error) sendMessageError {
	return SendMessageError(err)
}
