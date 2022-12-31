package meepo_core

import (
	"context"
	"crypto/ed25519"
	"net"
	"testing"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	meepo_testing "github.com/PeerXu/meepo/pkg/internal/testing"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_http "github.com/PeerXu/meepo/pkg/lib/rpc/http"
	"github.com/PeerXu/meepo/pkg/lib/stun"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	"github.com/PeerXu/meepo/pkg/meepo/tracker"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

func TestMeepo(t *testing.T) {
	ctx := context.Background()

	es, err := meepo_testing.NewEchoServer()
	assert.Nil(t, err)
	go es.Serve(ctx)        // nolint:errcheck
	defer es.Terminate(ctx) // nolint:errcheck

	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	assert.Nil(t, err)

	pubk, prik, _ := ed25519.GenerateKey(nil)
	signer := crypto_core.NewSigner(pubk, prik)
	cryptor := crypto_core.NewCryptor(pubk, prik, nil)
	addr, err := addr.FromBytesWithoutMagicCode(pubk)
	assert.Nil(t, err)

	var cfg webrtc.Configuration

	mp, err := NewMeepo(
		well_known_option.WithAddr(addr),
		well_known_option.WithLogger(logger),
		tracker_core.WithTrackers(),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithWebrtcConfiguration(cfg),
	)
	assert.Nil(t, err)
	defer mp.Close(ctx)

	tsp, err := mp.NewTransport(ctx, addr)
	assert.Nil(t, err)

	c, err := tsp.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("raw"))
	assert.Nil(t, err)

	err = c.WaitReady()
	assert.Nil(t, err)

	n, err := c.Conn().Write([]byte("hello, world"))
	assert.Nil(t, err)
	assert.Equal(t, 12, n)

	buf := make([]byte, n)
	n, err = c.Conn().Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, 12, n)
	assert.Equal(t, []byte("hello, world"), buf)
}

func TestMeepo2(t *testing.T) {
	ctx := context.Background()
	es, err := meepo_testing.NewEchoServer()
	require.Nil(t, err)
	go es.Serve(ctx)        // nolint:errcheck
	defer es.Terminate(ctx) // nolint:errcheck
	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	require.Nil(t, err)
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)
	defer lis.Close() // nolint:errcheck

	var mp1, mp2 meepo_interface.Meepo
	{
		pubk, prik, err := ed25519.GenerateKey(nil)
		require.Nil(t, err)
		signer := crypto_core.NewSigner(pubk, prik)
		cryptor := crypto_core.NewCryptor(pubk, prik, nil)
		addr, err := addr.FromBytesWithoutMagicCode(pubk)
		require.Nil(t, err)
		cfg := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{URLs: stun.STUNS},
			},
		}
		mp1, err = NewMeepo(
			well_known_option.WithAddr(addr),
			well_known_option.WithLogger(logger),
			tracker_core.WithTrackers(),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithWebrtcConfiguration(cfg),
			WithEnablePoof(false),
		)
		require.Nil(t, err)
		defer mp1.Close(ctx)

		srv, err := rpc.NewServer("http",
			rpc_core.WithHandler(mp1.AsTrackerdHandler()),
			well_known_option.WithLogger(logger),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithListener(lis),
		)
		require.Nil(t, err)
		go srv.Serve(ctx)        // nolint:errcheck
		defer srv.Terminate(ctx) // nolint:errcheck
	}

	{
		pubk, prik, err := ed25519.GenerateKey(nil)
		require.Nil(t, err)
		signer := crypto_core.NewSigner(pubk, prik)
		cryptor := crypto_core.NewCryptor(pubk, prik, nil)
		addr, err := addr.FromBytesWithoutMagicCode(pubk)
		require.Nil(t, err)
		cfg := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{URLs: stun.STUNS},
			},
		}
		caller, err := rpc.NewCaller("http",
			well_known_option.WithLogger(logger),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			rpc_http.WithBaseURL("http://"+lis.Addr().String()),
		)
		require.Nil(t, err)
		tk, err := tracker.NewTracker("rpc",
			well_known_option.WithAddr(mp1.Addr()),
			rpc_core.WithCaller(caller),
		)
		require.Nil(t, err)

		mp2, err = NewMeepo(
			well_known_option.WithAddr(addr),
			well_known_option.WithLogger(logger),
			tracker_core.WithTrackers(tk),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithWebrtcConfiguration(cfg),
		)
		require.Nil(t, err)
	}

	tsp, err := mp2.NewTransport(ctx, mp1.Addr())
	require.Nil(t, err)
	defer tsp.Close(ctx)

	err = tsp.WaitReady()
	require.Nil(t, err)

	c, err := tsp.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("raw"))
	require.Nil(t, err)

	err = c.WaitReady()
	require.Nil(t, err)

	conn := c.Conn()
	_, err = conn.Write([]byte("hello, world"))
	require.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	require.Nil(t, err)
	assert.Equal(t, []byte("hello, world"), buf[:n])
}

func SkipMeepo3(t *testing.T) {
	mpdCount := 1
	mpCount := 32
	ctx := context.Background()
	var tks []Tracker

	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	require.Nil(t, err)

	for m := 0; m < mpdCount; m++ {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		require.Nil(t, err)
		defer lis.Close() // nolint:errcheck
		pubk, prik, err := ed25519.GenerateKey(nil)
		require.Nil(t, err)
		signer := crypto_core.NewSigner(pubk, prik)
		cryptor := crypto_core.NewCryptor(pubk, prik, nil)
		addr, err := addr.FromBytesWithoutMagicCode(pubk)
		require.Nil(t, err)
		cfg := webrtc.Configuration{}
		mp, err := NewMeepo(
			well_known_option.WithAddr(addr),
			well_known_option.WithLogger(logger),
			tracker_core.WithTrackers(),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithWebrtcConfiguration(cfg),
			WithEnablePoof(false),
		)
		require.Nil(t, err)
		defer mp.Close(ctx) // nolint:staticcheck

		srv, err := rpc.NewServer("http",
			rpc_core.WithHandler(mp.AsTrackerdHandler()),
			well_known_option.WithLogger(logger),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithListener(lis),
		)
		require.Nil(t, err)
		go srv.Serve(ctx)        // nolint:errcheck
		defer srv.Terminate(ctx) // nolint:errcheck

		caller, err := rpc.NewCaller("http",
			well_known_option.WithLogger(logger),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			rpc_http.WithBaseURL("http://"+lis.Addr().String()),
		)
		require.Nil(t, err)
		tk, err := tracker.NewTracker("rpc",
			well_known_option.WithAddr(mp.Addr()),
			rpc_core.WithCaller(caller),
		)
		require.Nil(t, err)

		tks = append(tks, tk)
	}

	for m := 0; m < mpCount; m++ {
		pubk, prik, err := ed25519.GenerateKey(nil)
		require.Nil(t, err)
		signer := crypto_core.NewSigner(pubk, prik)
		cryptor := crypto_core.NewCryptor(pubk, prik, nil)
		addr, err := addr.FromBytesWithoutMagicCode(pubk)
		require.Nil(t, err)
		cfg := webrtc.Configuration{}
		mp, err := NewMeepo(
			well_known_option.WithAddr(addr),
			well_known_option.WithLogger(logger),
			tracker_core.WithTrackers(tks...),
			crypto_core.WithSigner(signer),
			crypto_core.WithCryptor(cryptor),
			marshaler.WithMarshaler(marshaler_json.Marshaler),
			marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			well_known_option.WithWebrtcConfiguration(cfg),
		)
		require.Nil(t, err)
		defer mp.Close(ctx) // nolint:staticcheck
		time.Sleep(31 * time.Millisecond)
	}

	time.Sleep(1000 * time.Second)
}
