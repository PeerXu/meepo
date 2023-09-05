package rpc_simple_http

import (
	"io"

	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	"github.com/gorilla/websocket"
)

type messageParser struct {
	unmarshaler marshaler_interface.Unmarshaler
}

func (p *messageParser) FromBytes(b []byte) (rpc_interface.Message, error) {
	return &message{b, p.unmarshaler}, nil
}

type messageMarshaler struct {
	marshaler marshaler_interface.Marshaler
	parser    rpc_interface.MessageParser
}

func (m *messageMarshaler) Marshal(v any) (rpc_interface.Message, error) {
	p, err := m.marshaler.Marshal(v)
	if err != nil {
		return nil, err
	}

	return m.parser.FromBytes(p)
}

type message struct {
	p           []byte
	unmarshaler marshaler_interface.Unmarshaler
}

func (m *message) Unmarshal(v any) error {
	return m.unmarshaler.Unmarshal(m.p, v)
}

func (m *message) ToBytes() ([]byte, error) {
	return m.p, nil
}

type wsStream struct {
	conn      *websocket.Conn
	marshaler marshaler_interface.Marshaler
	parser    rpc_interface.MessageParser
}

func (stm *wsStream) Marshaler() rpc_interface.MessageMarshaler {
	return &messageMarshaler{
		marshaler: stm.marshaler,
		parser:    stm.parser,
	}
}

func (stm *wsStream) RecvMessage() (rpc_interface.Message, error) {
	_, p, err := stm.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			return nil, io.EOF
		}

		return nil, err
	}

	return stm.parser.FromBytes(p)
}

func (stm *wsStream) SendMessage(m rpc_interface.Message) error {
	p, err := m.ToBytes()
	if err != nil {
		return err
	}

	return stm.conn.WriteMessage(websocket.BinaryMessage, p)
}

func (stm *wsStream) Close() error {
	return stm.conn.Close()
}
