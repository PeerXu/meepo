package rpc_interface

import (
	"io"
)

type MessageMarshaler interface {
	Marshal(v any) (Message, error)
}

type MessageParser interface {
	FromBytes([]byte) (Message, error)
}

type Message interface {
	Unmarshal(v any) error
	ToBytes() ([]byte, error)
}

type Stream interface {
	Marshaler() MessageMarshaler
	RecvMessage() (Message, error)
	SendMessage(Message) error
	io.Closer
}
