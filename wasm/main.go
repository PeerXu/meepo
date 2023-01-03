//go:build wasm && js

package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
	"syscall/js"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/acl"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/constant"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_http "github.com/PeerXu/meepo/pkg/lib/rpc/http"
	"github.com/PeerXu/meepo/pkg/lib/stun"
	mpo_webrtc "github.com/PeerXu/meepo/pkg/lib/webrtc"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_core "github.com/PeerXu/meepo/pkg/meepo/core"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	"github.com/PeerXu/meepo/pkg/meepo/tracker"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

type myAddr struct {
	network string
	address string
}

func (x *myAddr) Network() string { return x.network }
func (x *myAddr) String() string  { return x.address }

type conn struct {
	ch meepo_interface.Channel
}

func (c *conn) Read(b []byte) (n int, err error)   { return c.ch.Conn().Read(b) }
func (c *conn) Write(b []byte) (n int, err error)  { return c.ch.Conn().Write(b) }
func (c *conn) Close() error                       { return c.ch.Conn().Close() }
func (c *conn) LocalAddr() net.Addr                { return &myAddr{"meepo", "meepo"} }
func (c *conn) RemoteAddr() net.Addr               { return c.ch.SinkAddr() }
func (c *conn) SetDeadline(t time.Time) error      { return nil }
func (c *conn) SetReadDeadline(t time.Time) error  { return nil }
func (c *conn) SetWriteDeadline(t time.Time) error { return nil }

func NewMeepo(trackerdAddrStr string, trackerdURL string) (meepo_interface.Meepo, error) {
	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	if err != nil {
		return nil, err
	}

	pubk, prik, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	signer := crypto_core.NewSigner(pubk, prik)
	cryptor := crypto_core.NewCryptor(pubk, prik, nil)
	mpAddr, err := addr.FromBytesWithoutMagicCode(pubk)
	if err != nil {
		return nil, err
	}
	iceServers, err := mpo_webrtc.ParseICEServers(stun.STUNS)
	if err != nil {
		return nil, err
	}
	webrtcConfiguration := webrtc.Configuration{ICEServers: iceServers}

	tkAddr, err := addr.FromString(trackerdAddrStr)
	if err != nil {
		return nil, err
	}

	caller, err := rpc.NewCaller("http",
		well_known_option.WithLogger(logger),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		rpc_http.WithBaseURL(trackerdURL),
	)
	if err != nil {
		return nil, err
	}

	tk, err := tracker.NewTracker("rpc",
		well_known_option.WithAddr(tkAddr),
		rpc_core.WithCaller(caller),
	)
	tks := []tracker_core.Tracker{tk}

	// TODO: load from input
	acl_, err := acl.FromString(`- allow: "*"`)
	if err != nil {
		return nil, err
	}

	nmOpts := []meepo_core.NewMeepoOption{
		well_known_option.WithAddr(mpAddr),
		well_known_option.WithLogger(logger),
		tracker_core.WithTrackers(tks...),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithWebrtcConfiguration(webrtcConfiguration),
		acl.WithAcl(acl_),
		well_known_option.WithEnableMux(true),
		well_known_option.WithEnableKcp(false),

		well_known_option.WithMuxVer(constant.SMUX_VERSION),
		well_known_option.WithMuxBuf(constant.SMUX_BUFFER_SIZE),
		well_known_option.WithMuxStreamBuf(constant.SMUX_STREAM_BUFFER_SIZE),
		well_known_option.WithMuxNocomp(constant.SMUX_NOCOMP),
	}

	mp, err := meepo_core.NewMeepo(nmOpts...)
	if err != nil {
		logger.Fatal(err)
	}

	return mp, nil
}

func main() {
	c := make(chan struct{}, 0)

	var mp meepo_interface.Meepo
	js.Global().Get("document").Set("start", js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			preOutput := js.Global().Get("document").Call("getElementById", "output")
			preError := js.Global().Get("document").Call("getElementById", "error")

			preOutput.Set("innerHTML", "")
			preError.Set("innerHTML", "")

			inputTrackerdAddr := js.Global().Get("document").Call("getElementById", "trackerdAddr")
			inputTrackerdURL := js.Global().Get("document").Call("getElementById", "trackerdURL")
			trackerdAddr := inputTrackerdAddr.Get("value").String()
			trackerdURL := inputTrackerdURL.Get("value").String()

			var err error
			if mp, err = NewMeepo(trackerdAddr, trackerdURL); err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			preOutput.Set("innerHTML", "meepo started")
			js.Global().Get("document").Call("getElementById", "id").Set("innerHTML", mp.Addr().String())
			js.Global().Get("document").Call("getElementById", "startButton").Set("disabled", true)
		}()
		return false
	}))

	js.Global().Get("document").Set("listTransports", js.FuncOf(func(this js.Value, args []js.Value) any {
		preOutput := js.Global().Get("document").Call("getElementById", "output")
		preError := js.Global().Get("document").Call("getElementById", "error")

		preOutput.Set("innerHTML", "")
		preError.Set("innerHTML", "")

		ts, err := mp.ListTransports(context.Background())
		if err != nil {
			preError.Set("innerHTML", err.Error())
			return false
		}

		sort.Slice(ts, func(i, j int) bool { return ts[i].Addr().String() < ts[j].Addr().String() })

		var sb strings.Builder
		tb := tablewriter.NewWriter(&sb)
		tb.SetHeader([]string{"Addr", "State"})
		for _, t := range ts {
			tb.Append([]string{t.Addr().String(), t.State().String()})
		}
		tb.Render()

		preOutput.Set("innerHTML", sb.String())

		return false
	}))

	js.Global().Get("document").Set("newTransport", js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			ctx := context.Background()

			preOutput := js.Global().Get("document").Call("getElementById", "output")
			preError := js.Global().Get("document").Call("getElementById", "error")

			preOutput.Set("innerHTML", "")
			preError.Set("innerHTML", "")

			inputTarget := js.Global().Get("document").Call("getElementById", "target")
			targetStr := inputTarget.Get("value").String()
			target, err := addr.FromString(targetStr)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			t, err := mp.NewTransport(ctx, target)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			preOutput.Set("innerHTML", t.Addr().String())

			return
		}()
		return false
	}))

	js.Global().Get("document").Set("closeTransport", js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			ctx := context.Background()

			preOutput := js.Global().Get("document").Call("getElementById", "output")
			preError := js.Global().Get("document").Call("getElementById", "error")

			preOutput.Set("innerHTML", "")
			preError.Set("innerHTML", "")

			inputTarget := js.Global().Get("document").Call("getElementById", "target")
			targetStr := inputTarget.Get("value").String()
			target, err := addr.FromString(targetStr)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			t, err := mp.GetTransport(ctx, target)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			err = t.Close(ctx)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			preOutput.Set("innerHTML", fmt.Sprintf("%s closed", target.String()))

			return
		}()
		return false
	}))

	js.Global().Get("document").Set("whoami", js.FuncOf(func(this js.Value, args []js.Value) any {
		preOutput := js.Global().Get("document").Call("getElementById", "output")
		preError := js.Global().Get("document").Call("getElementById", "error")

		preOutput.Set("innerHTML", "")
		preError.Set("innerHTML", "")

		preOutput.Set("innerHTML", mp.Addr().String())

		return false
	}))

	js.Global().Get("document").Set("ping", js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			ctx := context.Background()

			preOutput := js.Global().Get("document").Call("getElementById", "output")
			preError := js.Global().Get("document").Call("getElementById", "error")

			preOutput.Set("innerHTML", "")
			preError.Set("innerHTML", "")

			inputTarget := js.Global().Get("document").Call("getElementById", "target")
			targetStr := inputTarget.Get("value").String()
			target, err := addr.FromString(targetStr)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			t, err := mp.GetTransport(ctx, target)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			var res meepo_core.PingResponse
			if err = t.Call(ctx, "ping", &meepo_core.PingRequest{Nonce: 0}, &res); err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			preOutput.Set("innerHTML", fmt.Sprintf("ping: %v", time.Now()))
		}()

		return false
	}))

	js.Global().Get("document").Set("doRequest", js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			ctx := context.Background()

			preOutput := js.Global().Get("document").Call("getElementById", "output")
			preError := js.Global().Get("document").Call("getElementById", "error")

			preOutput.Set("innerHTML", "")
			preError.Set("innerHTML", "")

			inputTarget := js.Global().Get("document").Call("getElementById", "target")
			inputMode := js.Global().Get("document").Call("getElementById", "mode")
			inputMethod := js.Global().Get("document").Call("getElementById", "method")
			inputURL := js.Global().Get("document").Call("getElementById", "url")
			inputHeaders := js.Global().Get("document").Call("getElementById", "headers")
			inputBody := js.Global().Get("document").Call("getElementById", "body")

			targetStr := inputTarget.Get("value").String()
			target, err := addr.FromString(targetStr)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			mode := inputMode.Get("value").String()
			method := inputMethod.Get("value").String()
			url := inputURL.Get("value").String()
			headersStr := inputHeaders.Get("value").String()
			headers := make(map[string]string)
			if err = json.Unmarshal([]byte(headersStr), &headers); err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}
			body := inputBody.Get("value").String()

			httpTransport := &http.Transport{
				Dial: func(network, address string) (net.Conn, error) {
					t, err := mp.GetTransport(ctx, target)
					if err != nil {
						return nil, err
					}

					ch, err := t.NewChannel(ctx, network, address, well_known_option.WithMode(mode))
					if err != nil {
						return nil, err
					}

					if err = ch.WaitReady(); err != nil {
						return nil, err
					}

					return &conn{ch}, nil
				},
			}
			cli := http.Client{Transport: httpTransport}
			req, err := http.NewRequest(method, url, strings.NewReader(body))
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			res, err := cli.Do(req)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}
			defer res.Body.Close()

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				preError.Set("innerHTML", err.Error())
				return
			}

			preOutput.Set("innerHTML", string(resBody))
		}()

		return false
	}))

	<-c
}
