package rpc_simple_http

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/rand"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (c *SimpleHttpCaller) CallStream(ctx context.Context, method string, opts ...rpc_interface.CallStreamOption) (rpc_interface.Stream, error) {
	session := rand.DefaultStringGenerator.Generate(8)

	logger := c.GetLogger().WithFields(logging.Fields{
		"#method": "CallStream",
		"method":  method,
		"session": session,
	})

	urlStr := c.JoinPath("/v1/actions/simpleDoStream")
	urlStr = fmt.Sprintf("%s?method=%s&session=%s", urlStr, method, session)

	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to dial")
		return nil, err
	}

	stm, err := c.wrapWsConnToStream(conn)
	if err != nil {
		logger.WithError(err).Debugf("failed to wrap websocket connection to stream")
		return nil, err
	}

	return stm, nil
}

func (c *SimpleHttpCaller) wrapWsConnToStream(conn *websocket.Conn) (rpc_interface.Stream, error) {
	return &wsStream{
		conn:      conn,
		marshaler: c.marshaler,
		parser:    &messageParser{c.unmarshaler},
	}, nil
}
