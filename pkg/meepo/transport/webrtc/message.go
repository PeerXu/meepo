package transport_webrtc

import (
	"fmt"
)

type Message struct {
	Session string
	Scope   string
	Method  string
	Data    []byte
	Error   string
}

func NewMessage(sess string, scope, method string, data []byte, errStr string) Message {
	return Message{
		Session: sess,
		Scope:   scope,
		Method:  method,
		Data:    data,
		Error:   errStr,
	}
}

func (t *WebrtcTransport) newRequest(scope, method string, data []byte) Message {
	return NewMessage(t.newRequestSession(), scope, method, data, "")
}

func (t *WebrtcTransport) NewResponse(req Message, data []byte) Message {
	return NewMessage(t.parseResponseSession(req.Session), req.Scope, req.Method, data, "")
}

func (t *WebrtcTransport) NewErrorResponse(req Message, err error) Message {
	return NewMessage(t.parseResponseSession(req.Session), req.Scope, req.Method, nil, err.Error())
}

func (t *WebrtcTransport) newRequestSession() string {
	return fmt.Sprintf("%08x", uint32(t.randSrc.Int63()|0x1))
}

func (t *WebrtcTransport) parseResponseSession(reqSess string) string {
	return parseResponseSession(reqSess)
}
