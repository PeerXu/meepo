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
	"honnef.co/go/js/dom/v2"

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

type conn struct{ meepo_interface.Channel }

func (c *conn) Read(b []byte) (n int, err error)   { return c.Conn().Read(b) }
func (c *conn) Write(b []byte) (n int, err error)  { return c.Conn().Write(b) }
func (c *conn) Close() error                       { return c.Conn().Close() }
func (c *conn) LocalAddr() net.Addr                { return &net.UnixAddr{Net: "meepo", Name: "meepo"} }
func (c *conn) RemoteAddr() net.Addr               { return c.SinkAddr() }
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

	acl_, err := acl.FromString(`- block: "*"`)
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
	doc := dom.GetWindow().Document()
	preOutput := doc.GetElementByID("output").(*dom.HTMLPreElement)
	preError := doc.GetElementByID("error").(*dom.HTMLPreElement)
	clearScreen := func() {
		preOutput.SetInnerHTML("")
		preError.SetInnerHTML("")
	}
	printError := func(err error) {
		clearScreen()
		preError.SetInnerHTML(err.Error())
	}
	printMessage := func(msg string) {
		clearScreen()
		preOutput.SetInnerHTML(msg)
	}

	simpleNewJsFunc := func(name string, fn func()) {
		doc.Underlying().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
			go fn()
			return false
		}))
	}

	simpleNewJsFunc("start", func() {
		var err error

		inputTrackerdAddr := doc.GetElementByID("trackerdAddr").(*dom.HTMLInputElement)
		inputTrackerdURL := doc.GetElementByID("trackerdURL").(*dom.HTMLInputElement)
		trackerdAddr := inputTrackerdAddr.Value()
		trackerdURL := inputTrackerdURL.Value()

		if mp, err = NewMeepo(trackerdAddr, trackerdURL); err != nil {
			printError(err)
			return
		}

		printMessage("meepo started")

		doc.GetElementByID("id").(*dom.HTMLSpanElement).SetInnerHTML(mp.Addr().String())
		startButton := doc.GetElementByID("startButton").(*dom.HTMLButtonElement)
		startButton.SetInnerHTML("running")
		startButton.Set("disabled", true)

		for _, e := range doc.QuerySelectorAll(".enableAfterRunning") {
			e.(*dom.HTMLButtonElement).Set("disabled", false)
		}
	})

	simpleNewJsFunc("listTransports", func() {
		ts, err := mp.ListTransports(context.Background())
		if err != nil {
			printError(err)
			return
		}

		sort.Slice(ts, func(i, j int) bool { return ts[i].Addr().String() < ts[j].Addr().String() })

		var sb strings.Builder
		tb := tablewriter.NewWriter(&sb)
		tb.SetHeader([]string{"Addr", "State"})
		for _, t := range ts {
			tb.Append([]string{t.Addr().String(), t.State().String()})
		}
		tb.Render()

		printMessage(sb.String())
	})

	simpleNewJsFunc("newTransport", func() {
		ctx := context.Background()

		inputTarget := doc.GetElementByID("target").(*dom.HTMLInputElement)
		targetStr := inputTarget.Value()
		target, err := addr.FromString(targetStr)
		if err != nil {
			printError(err)
			return
		}

		t, err := mp.NewTransport(ctx, target)
		if err != nil {
			printError(err)
			return
		}

		printMessage(fmt.Sprintf("%s created", t.Addr().String()))
		return
	})

	simpleNewJsFunc("closeTransport", func() {
		ctx := context.Background()

		inputTarget := doc.GetElementByID("target").(*dom.HTMLInputElement)
		targetStr := inputTarget.Value()
		target, err := addr.FromString(targetStr)
		if err != nil {
			printError(err)
			return
		}

		t, err := mp.GetTransport(ctx, target)
		if err != nil {
			printError(err)
			return
		}

		err = t.Close(ctx)
		if err != nil {
			printError(err)
			return
		}

		printMessage(fmt.Sprintf("%s closed", target.String()))
		return
	})

	simpleNewJsFunc("whoami", func() { printMessage(mp.Addr().String()) })

	simpleNewJsFunc("ping", func() {
		ctx := context.Background()

		inputTarget := doc.GetElementByID("target").(*dom.HTMLInputElement)
		targetStr := inputTarget.Value()
		target, err := addr.FromString(targetStr)
		if err != nil {
			printError(err)
			return
		}

		t, err := mp.GetTransport(ctx, target)
		if err != nil {
			printError(err)
			return
		}

		var res meepo_core.PingResponse
		if err = t.Call(ctx, "ping", &meepo_core.PingRequest{Nonce: 0}, &res); err != nil {
			printError(err)
			return
		}

		printMessage(fmt.Sprintf("ping: %v", time.Now()))
		return
	})

	simpleNewJsFunc("doRequest", func() {
		ctx := context.Background()

		inputTarget := doc.GetElementByID("target").(*dom.HTMLInputElement)
		targetStr := inputTarget.Value()
		target, err := addr.FromString(targetStr)
		if err != nil {
			preError.Set("innerHTML", err.Error())
			return
		}

		inputMode := doc.GetElementByID("mode").(*dom.HTMLInputElement)
		mode := inputMode.Value()

		inputMethod := doc.GetElementByID("method").(*dom.HTMLInputElement)
		method := inputMethod.Value()

		inputURL := doc.GetElementByID("url").(*dom.HTMLInputElement)
		url := inputURL.Value()

		inputHeaders := doc.GetElementByID("headers").(*dom.HTMLInputElement)
		headersStr := inputHeaders.Value()
		headers := make(map[string]string)
		if err = json.Unmarshal([]byte(headersStr), &headers); err != nil {
			preError.Set("innerHTML", err.Error())
			return
		}

		inputBody := doc.GetElementByID("body").(*dom.HTMLInputElement)
		body := inputBody.Value()

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
			printError(err)
			return
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		res, err := cli.Do(req)
		if err != nil {
			printError(err)
			return
		}
		defer res.Body.Close()

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			printError(err)
			return
		}

		printMessage(string(resBody))
		return
	})

	<-c
}
