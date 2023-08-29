package rpc_simple_http

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (s *SimpleHttpServer) HttpSimpleDoStream(c *gin.Context) {
	logger := s.GetLogger().WithField("#method", "HttpSimpleDoStream")

	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.WithError(err).Errorf("failed to upgrade connection")
		return
	}
	defer conn.Close()

	stm, err := s.wrapWsToStream(conn)
	if err != nil {
		logger.WithError(err).Errorf("failed to wrap websocket connection to stream")
		return
	}

	m, err := stm.RecvMessage()
	if err != nil {
		logger.WithError(err).Errorf("failed to receive setup message")
		return
	}

	var req InitRequest
	if err = m.Unmarshal(&req); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal message to initial request")
		return
	}
	logger = logger.WithFields(logging.Fields{
		"method":  req.Method,
		"session": req.Session,
	})

	m, err = stm.Marshaler().Marshal(&InitResponse{Session: req.Session})
	if err != nil {
		logger.WithError(err).Errorf("failed to marshal initial response to message")
		return
	}

	if stm.SendMessage(m); err != nil {
		logger.WithError(err).Errorf("failed to send initial response message")
		return
	}

	logger.Debugf("simple do stream")

	ctx := context.WithValue(s.context(), CONTEXT_SESSION, req.Session)
	if err = s.handler.DoStream(ctx, req.Method, stm); err != nil {
		logger.WithError(err).Errorf("failed to do stream")
		return
	}

	logger.Debugf("simple do stream done")
}

type InitRequest struct {
	Session string `json:"session"`
	Method  string `json:"method"`
}

type InitResponse struct {
	Session string `json:"session"`
}

// TODO: avoid to alloc too many objects
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

func (s *SimpleHttpServer) wrapWsToStream(conn *websocket.Conn) (rpc_interface.Stream, error) {
	return &wsStream{
		conn:      conn,
		marshaler: s.marshaler,
		parser:    &messageParser{s.unmarshaler},
	}, nil
}
