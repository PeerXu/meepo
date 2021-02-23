package meepo

import (
	"fmt"
)

var (
	GatherTimeoutError         = fmt.Errorf("Gather timeout")
	TransportNotExistError     = fmt.Errorf("Transport not exists")
	TransportExistError        = fmt.Errorf("Transport already exists")
	TeleportationNotExistError = fmt.Errorf("Teleportation not exists")
	UnexpectedMessageError     = fmt.Errorf("Unexpected message")
	WaitResponseTimeoutError   = fmt.Errorf("Wait response timeout")
	NotListenableAddressError  = fmt.Errorf("Not listenable address")
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

func UnsupportedRequestHandlerError(method string) error {
	return fmt.Errorf("Unsupported method: %s", method)
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
