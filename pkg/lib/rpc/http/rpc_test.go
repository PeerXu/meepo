package rpc_http

import (
	"context"
	"crypto/ed25519"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_default "github.com/PeerXu/meepo/pkg/lib/rpc/default"
)

func TestHttp(t *testing.T) {
	pubk, prik, _ := ed25519.GenerateKey(nil)
	logger, _ := logging.NewLogger()
	signer := crypto_core.NewSigner(pubk, prik)
	cryptor := crypto_core.NewCryptor(pubk, prik, nil)

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

	s, err := NewHttpServer(
		rpc_core.WithHandler(handler),
		well_known_option.WithLogger(logger),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithListener(lis),
	)
	require.Nil(t, err)
	go s.Serve(context.Background())
	defer s.Terminate(context.Background()) // nolint:errcheck

	caller, err := NewHttpCaller(
		well_known_option.WithLogger(logger),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		WithBaseURL("http://"+lis.Addr().String()),
	)
	require.Nil(t, err)

	var res string
	err = caller.Call(context.Background(), "echo", map[string]string{"text": "hello, world!"}, &res, well_known_option.WithDestination(pubk))
	assert.Nil(t, err)
	assert.Equal(t, "hello, world!", res)
}
