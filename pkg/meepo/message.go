package meepo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

type MessageGetter interface {
	GetMessage() Message
}

type Message struct {
	PeerID  string `json:"_peerID"`
	Type    string `json:"_type"`
	Session int32  `json:"_session"`
	Method  string `json:"_method"`
	Error   string `json:"_error,omitempty"`
}

func (m Message) GetMessage() Message {
	return m
}

func IsMessage(m *Message) bool {
	return m.PeerID != "" &&
		(m.Type == "request" || m.Type == "response")
}

func InvertMessage(m Message, id string) Message {
	var typ string
	switch m.Type {
	case "request":
		typ = "response"
	case "response":
		typ = "request"
	}

	return Message{
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
		if !IsMessage(&m) {
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
