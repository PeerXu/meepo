package meepo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

const (
	MESSAGE_TYPE_REQUEST           = "request"
	MESSAGE_TYPE_BROADCAST_REQUEST = "broadcast_request"
	MESSAGE_TYPE_RESPONSE          = "response"
)

type MessageGetter interface {
	GetMessage() *Message
}

type Message struct {
	PeerID  string `json:"_peerID"`
	Type    string `json:"_type"`
	Session int32  `json:"_session"`
	Method  string `json:"_method"`
	Error   string `json:"_error,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("#<Message: PeerID: %v, Type: %v, Session: %v, Method: %v, Error: %v>",
		m.PeerID, m.Type, m.Session, m.Method, m.Error)
}

func (m *Message) Identifier() string {
	return fmt.Sprintf("%v.%v", m.Method, m.Session)
}

func IsMessage(m *Message) bool {
	return m.PeerID != "" &&
		(m.Type == MESSAGE_TYPE_REQUEST ||
			m.Type == MESSAGE_TYPE_RESPONSE ||
			m.Type == MESSAGE_TYPE_BROADCAST_REQUEST)
}

func InvertMessage(m *Message, id string) *Message {
	var typ string
	switch m.Type {
	case MESSAGE_TYPE_REQUEST:
		fallthrough
	case MESSAGE_TYPE_BROADCAST_REQUEST:
		typ = MESSAGE_TYPE_RESPONSE
	case MESSAGE_TYPE_RESPONSE:
		typ = MESSAGE_TYPE_REQUEST
	}

	return &Message{
		PeerID:  id,
		Type:    typ,
		Session: m.Session,
		Method:  m.Method,
	}
}

type decodeMessageFunc func([]byte) (interface{}, error)

var decodeMessageFuncs sync.Map

func messageIdentifier(m Message) string {
	return joinTypeMethodIdentifier(m.Type, m.Method)
}

func joinTypeMethodIdentifier(typ, meth string) string {
	return fmt.Sprintf("%s.%s", typ, meth)
}

func registerDecodeMessage(typ, meth string, fn decodeMessageFunc) {
	decodeMessageFuncs.Store(joinTypeMethodIdentifier(typ, meth), fn)
}

func registerDecodeMessageHelper(typ, meth string, fn func() interface{}) {
	registerDecodeMessage(typ, meth, func(buf []byte) (interface{}, error) {
		var err error

		in := fn()

		if err = json.NewDecoder(bytes.NewReader(buf)).Decode(in); err != nil {
			return nil, err
		}

		msgGetter, ok := in.(MessageGetter)
		if !ok {
			return nil, UnexpectedMessageError
		}

		m := msgGetter.GetMessage()
		if !IsMessage(m) {
			return nil, UnexpectedMessageError
		}

		return in, nil
	})
}

func DecodeMessage(buf []byte) (interface{}, error) {
	var err error
	var m Message

	if err = json.NewDecoder(bytes.NewReader(buf)).Decode(&m); err != nil {
		return nil, err
	}

	mid := messageIdentifier(m)
	fn, ok := decodeMessageFuncs.Load(mid)
	if !ok {
		return nil, UnsupportedMessageDecodeDriverError(mid)
	}

	return fn.(decodeMessageFunc)(buf)
}
