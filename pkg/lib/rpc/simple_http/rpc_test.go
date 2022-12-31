package rpc_simple_http

import (
	"context"
	"net"
	"testing"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_default "github.com/PeerXu/meepo/pkg/lib/rpc/default"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleHttp(t *testing.T) {
	logger, _ := logging.NewLogger(logging.WithLevel("trace"))
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)
	defer lis.Close()

	handler, err := rpc_default.NewDefaultHandler()
	require.Nil(t, err)
	handler.Handle("echo", func(ctx context.Context, in []byte) ([]byte, error) {
		req := map[string]string{}
		err := marshaler.Unmarshal(ctx, in, &req)
		require.Nil(t, err)
		return marshaler.Marshal(ctx, req["text"])
	})

	s, err := NewSimpleHttpServer(
		rpc_core.WithHandler(handler),
		well_known_option.WithLogger(logger),
		well_known_option.WithListener(lis),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
	)
	require.Nil(t, err)
	go s.Serve(context.Background())
	defer s.Terminate(context.Background()) // nolint:errcheck

	caller, err := NewSimpleHttpCaller(
		well_known_option.WithLogger(logger),
		WithBaseURL("http://"+lis.Addr().String()),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
	)
	require.Nil(t, err)

	var res string
	err = caller.Call(context.Background(), "echo", map[string]string{"text": "hello, world!"}, &res)
	assert.Nil(t, err)
	assert.Equal(t, "hello, world!", res)
}
