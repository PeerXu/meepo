package rpc_simple_http

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (s *SimpleHttpServer) HttpSimpleDoStream(c *gin.Context) {
	method := c.Query("method")

	logger := s.GetLogger().WithFields(logging.Fields{
		"#method": "HttpSimpleDoStream",
		"method":  method,
	})

	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.WithError(err).Errorf("failed to upgrade connection")
		return
	}
	defer conn.Close()

	stm, err := s.wrapWsConnToStream(conn)
	if err != nil {
		logger.WithError(err).Errorf("failed to wrap websocket connection to stream")
		return
	}

	logger.Debugf("simple do stream")

	if err = s.handler.DoStream(s.context(), method, stm); err != nil {
		logger.WithError(err).Errorf("failed to do stream")
		return
	}

	logger.Debugf("simple do stream done")
}

func (s *SimpleHttpServer) wrapWsConnToStream(conn *websocket.Conn) (rpc_interface.Stream, error) {
	return &wsStream{
		conn:      conn,
		marshaler: s.marshaler,
		parser:    &messageParser{s.unmarshaler},
	}, nil
}
