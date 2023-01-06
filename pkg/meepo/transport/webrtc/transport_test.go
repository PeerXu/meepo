package transport_webrtc

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"testing"

	"github.com/pion/webrtc/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	"github.com/PeerXu/meepo/pkg/lib/stun"
	meepo_testing "github.com/PeerXu/meepo/pkg/lib/testing"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func TestRawTransport(t *testing.T) {
	ctx := context.Background()
	es, err := meepo_testing.NewEchoServer()
	require.Nil(t, err)
	go es.Serve(ctx)        // nolint
	defer es.Terminate(ctx) // nolint

	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	require.Nil(t, err)

	var se webrtc.SettingEngine
	se.DetachDataChannels()
	api := webrtc.NewAPI(webrtc.WithSettingEngine(se))
	cfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: stun.STUNS,
			},
		},
	}

	var offer, answer *webrtc.SessionDescription
	var waitGather, waitGatherDone = make(chan struct{}), make(chan struct{})
	gather, gatherDone := func() (GatherFunc, GatherDoneFunc) {
		return func(_offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
				offer = &_offer
				close(waitGather)
				<-waitGatherDone
				return *answer, nil
			}, func(_answer webrtc.SessionDescription, er error) {
				require.Nil(t, er)
				answer = &_answer
				close(waitGatherDone)
			}
	}()

	addrSrcBuf, _, _ := ed25519.GenerateKey(nil)
	addrSrc := addr.Must(addr.FromBytesWithoutMagicCode(addrSrcBuf))
	addrSinkBuf, _, _ := ed25519.GenerateKey(nil)
	for i := 1; i < len(addrSinkBuf); i++ {
		addrSinkBuf[i] = 0xff
	}
	addrSink := addr.Must(addr.FromBytesWithoutMagicCode(addrSinkBuf))

	pcSrc, err := api.NewPeerConnection(cfg)
	require.Nil(t, err)

	pcSink, err := api.NewPeerConnection(cfg)
	require.Nil(t, err)

	tSrc, err := NewWebrtcSourceTransport(
		well_known_option.WithLogger(logger),
		well_known_option.WithAddr(addrSrc),
		WithGatherFunc(gather),
		well_known_option.WithPeerConnection(pcSrc),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		transport_core.WithBeforeNewChannelHook(func(t meepo_interface.Transport, network string, address string) error { return nil }),
		transport_core.WithOnTransportCloseFunc(func(meepo_interface.Transport) error { return nil }),
		transport_core.WithOnTransportReadyFunc(func(meepo_interface.Transport) error { return nil }),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithEnableMux(false),
		well_known_option.WithEnableKcp(false),
	)
	defer tSrc.Close(ctx) // nolint:staticcheck
	require.Nil(t, err)
	<-waitGather

	tSink, err := NewWebrtcSinkTransport(
		well_known_option.WithLogger(logger),
		well_known_option.WithAddr(addrSink),
		WithOffer(*offer),
		WithGatherDoneFunc(gatherDone),
		well_known_option.WithPeerConnection(pcSink),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		transport_core.WithBeforeNewChannelHook(func(t meepo_interface.Transport, network string, address string) error { return nil }),
		transport_core.WithOnTransportCloseFunc(func(meepo_interface.Transport) error { return nil }),
		transport_core.WithOnTransportReadyFunc(func(meepo_interface.Transport) error { return nil }),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithEnableMux(false),
		well_known_option.WithEnableKcp(false),
	)
	require.Nil(t, err)
	defer tSink.Close(ctx) // nolint:staticcheck

	<-waitGatherDone

	tSrc.WaitReady()  // nolint:errcheck
	tSink.WaitReady() // nolint:errcheck

	c1, err := tSrc.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("raw"))
	require.Nil(t, err)
	defer c1.Close(ctx) // nolint:staticcheck

	err = c1.WaitReady()
	require.Nil(t, err)

	conn1 := c1.Conn()
	_, err = conn1.Write([]byte("hello, world!"))
	require.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := conn1.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello, world!"), buf[:n])

	c2, err := tSrc.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("raw"))
	require.Nil(t, err)
	defer c2.Close(ctx) // nolint:staticcheck

	err = c2.WaitReady()
	require.Nil(t, err)

	conn2 := c2.Conn()
	_, err = conn2.Write([]byte("hello, world!"))
	require.Nil(t, err)

	buf = make([]byte, 1024)
	n, err = conn2.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello, world!"), buf[:n])
}

func TestMuxTransport(t *testing.T) {
	ctx := context.Background()
	es, err := meepo_testing.NewEchoServer()
	require.Nil(t, err)
	go es.Serve(ctx)        // nolint
	defer es.Terminate(ctx) // nolint

	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	require.Nil(t, err)

	var se webrtc.SettingEngine
	se.DetachDataChannels()
	api := webrtc.NewAPI(webrtc.WithSettingEngine(se))
	cfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: stun.STUNS,
			},
		},
	}

	var offer, answer *webrtc.SessionDescription
	var waitGather, waitGatherDone = make(chan struct{}), make(chan struct{})
	gather, gatherDone := func() (GatherFunc, GatherDoneFunc) {
		return func(_offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
				offer = &_offer
				close(waitGather)
				<-waitGatherDone
				return *answer, nil
			}, func(_answer webrtc.SessionDescription, er error) {
				require.Nil(t, er)
				answer = &_answer
				close(waitGatherDone)
			}
	}()

	muxLabel := fmt.Sprintf("mux#%016x", rand.Int63())
	addrSrcBuf, _, _ := ed25519.GenerateKey(nil)
	addrSrc := addr.Must(addr.FromBytesWithoutMagicCode(addrSrcBuf))
	addrSinkBuf, _, _ := ed25519.GenerateKey(nil)
	for i := 1; i < len(addrSinkBuf); i++ {
		addrSinkBuf[i] = 0xff
	}
	addrSink := addr.Must(addr.FromBytesWithoutMagicCode(addrSinkBuf))

	pcSrc, err := api.NewPeerConnection(cfg)
	require.Nil(t, err)

	pcSink, err := api.NewPeerConnection(cfg)
	require.Nil(t, err)

	tSrc, err := NewWebrtcSourceTransport(
		well_known_option.WithLogger(logger),
		well_known_option.WithAddr(addrSrc),
		WithGatherFunc(gather),
		well_known_option.WithPeerConnection(pcSrc),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		transport_core.WithBeforeNewChannelHook(func(t meepo_interface.Transport, network string, address string) error { return nil }),
		transport_core.WithOnTransportCloseFunc(func(meepo_interface.Transport) error { return nil }),
		transport_core.WithOnTransportReadyFunc(func(meepo_interface.Transport) error { return nil }),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithEnableMux(true),
		WithMuxLabel(muxLabel),
		well_known_option.WithEnableKcp(false),
	)
	require.Nil(t, err)
	defer tSrc.Close(ctx) // nolint:staticcheck

	<-waitGather

	tSink, err := NewWebrtcSinkTransport(
		well_known_option.WithLogger(logger),
		well_known_option.WithAddr(addrSink),
		WithOffer(*offer),
		WithGatherDoneFunc(gatherDone),
		well_known_option.WithPeerConnection(pcSink),
		dialer.WithDialer(dialer.GetGlobalDialer()),
		transport_core.WithBeforeNewChannelHook(func(t meepo_interface.Transport, network string, address string) error { return nil }),
		transport_core.WithOnTransportCloseFunc(func(meepo_interface.Transport) error { return nil }),
		transport_core.WithOnTransportReadyFunc(func(meepo_interface.Transport) error { return nil }),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithEnableMux(true),
		WithMuxLabel(muxLabel),
		well_known_option.WithEnableKcp(false),
	)
	require.Nil(t, err)
	defer tSink.Close(ctx) // nolint:staticcheck

	<-waitGatherDone

	tSrc.WaitReady()  // nolint:errcheck
	tSink.WaitReady() // nolint:errcheck

	c1, err := tSrc.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("mux"))
	require.Nil(t, err)
	defer c1.Close(ctx) // nolint:staticcheck

	err = c1.WaitReady()
	require.Nil(t, err)

	conn1 := c1.Conn()
	_, err = conn1.Write([]byte("hello, world!"))
	require.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := conn1.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello, world!"), buf[:n])

	c2, err := tSrc.NewChannel(ctx, es.Listener.Addr().Network(), es.Listener.Addr().String(), well_known_option.WithMode("raw"))
	require.Nil(t, err)
	defer c2.Close(ctx) // nolint:staticcheck

	err = c2.WaitReady()
	require.Nil(t, err)

	conn2 := c2.Conn()
	_, err = conn2.Write([]byte("hello, world!"))
	require.Nil(t, err)

	buf = make([]byte, 1024)
	n, err = conn2.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello, world!"), buf[:n])
}
