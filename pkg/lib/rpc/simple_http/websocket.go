package rpc_simple_http

import (
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

var (
	WebsocketUpgraderReadBufferSize  = 4 * 1024
	WebsocketUpgraderWriteBufferSize = 4 * 1024
	WebsocketUpgraderEnableWritePool = false

	upgrader *websocket.Upgrader
)

func init() {
	var err error
	rdbufsz := WebsocketUpgraderReadBufferSize
	if s := os.Getenv("MPO_EXPERIMENTAL_WEBSOCKET_UPGRADER_READ_BUFFER_SIZE"); s != "" {
		if rdbufsz, err = strconv.Atoi(s); err != nil {
			panic(err)
		}
	}
	wrbufsz := WebsocketUpgraderWriteBufferSize
	if s := os.Getenv("MPO_EXPERIMENTAL_WEBSOCKET_UPGRADER_WRITE_BUFFER_SIZE"); s != "" {
		if wrbufsz, err = strconv.Atoi(s); err != nil {
			panic(err)
		}
	}
	var pool websocket.BufferPool
	if s := os.Getenv("MPO_EXPERIMENTAL_WEBSOCKET_UPGRADER_ENABLE_WRITE_POOL"); s != "" {
		if enable := cast.ToBool(s); enable {
			pool = &sync.Pool{}
		}

	}

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  rdbufsz,
		WriteBufferSize: wrbufsz,
		WriteBufferPool: pool,
	}
}
