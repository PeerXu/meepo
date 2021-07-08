package meepo

import (
	"fmt"
)

var (
	ErrInvalidPeerID            = fmt.Errorf("invalid peer id")
	ErrTransportExist           = fmt.Errorf("transport exists")
	ErrTeleportationNotExist    = fmt.Errorf("teleportation not exists")
	ErrWaitResponseTimeout      = fmt.Errorf("wait response timeout")
	ErrNotWirable               = fmt.Errorf("not wirable")
	ErrUnsupportedSocks5Command = fmt.Errorf("unsupported socks5 command")
	ErrNetworkUnreachable       = fmt.Errorf("network unreachable")
	ErrUnsupportedNetworkType   = fmt.Errorf("unsupported network type")
	ErrNotBroadcastPacket       = fmt.Errorf("not braoadcast packet")
	ErrOutOfEdge                = fmt.Errorf("out of edge")
	ErrTransportNotExist        = fmt.Errorf("transport not exists")
	ErrUnexpectedMessage        = fmt.Errorf("unexpected message")
	ErrUnsupportedMethod        = fmt.Errorf("unsupported method")
	ErrUnexpectedType           = fmt.Errorf("unexpected type")
	ErrNotFound                 = fmt.Errorf("not found")
	ErrUnauthenticated          = fmt.Errorf("unauthenticated")
	ErrUnauthorized             = fmt.Errorf("unauthorized")
	ErrIncorrectSignature       = fmt.Errorf("incorrect signature")
	ErrIncorrectPassword        = fmt.Errorf("incorrect password")
	ErrAclNotAllowed            = fmt.Errorf("acl: not allowed")
	ErrInvalidAclPolicyString   = fmt.Errorf("invalid acl policy string")
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

type errSendPacket struct {
	err error
}

func (t errSendPacket) Error() string {
	return t.err.Error()
}

func ErrSendPacket(err error) errSendPacket {
	return ErrSendPacket(err)
}
