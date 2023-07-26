//go:build wasm && js

package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall/js"
	"time"

	"github.com/pion/webrtc/v3"
	"golang.org/x/crypto/ssh"
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

type myAddr struct {
	network string
	address string
}

func (x *myAddr) Network() string { return x.network }
func (x *myAddr) String() string  { return x.address }

type myConn struct{ meepo_interface.Channel }

func (c *myConn) Read(b []byte) (n int, err error)   { return c.Conn().Read(b) }
func (c *myConn) Write(b []byte) (n int, err error)  { return c.Conn().Write(b) }
func (c *myConn) Close() error                       { return c.Conn().Close() }
func (c *myConn) LocalAddr() net.Addr                { return &myAddr{"meepo", "meepo"} }
func (c *myConn) RemoteAddr() net.Addr               { return c.SinkAddr() }
func (c *myConn) SetDeadline(t time.Time) error      { return nil }
func (c *myConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *myConn) SetWriteDeadline(t time.Time) error { return nil }

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
	q := make(chan struct{}, 0)

	document := dom.GetWindow().Document()
	simpleNewFunction := func(name string, async bool, fn func(js.Value, []js.Value) any) {
		document.Underlying().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
			if async {
				go fn(this, args)
				return false
			}
			return fn(this, args)
		}))
	}
	getInputValue := func(id string) string {
		return document.GetElementByID(id).(*dom.HTMLInputElement).Value()
	}
	setElementInnerHTML := func(id string, body string) {
		document.GetElementByID(id).SetInnerHTML(body)
	}

	doitButton := document.GetElementByID("doit").(*dom.HTMLButtonElement)

	trackerdAddr := getInputValue("trackerdAddr")
	trackerdURL := getInputValue("trackerdURL")
	channelMode := "mux"

	var mp meepo_interface.Meepo
	var err error
	var sshSession *ssh.Session
	var sshSessionStdinChannel io.WriteCloser
	var sshSessionTarget string
	var sshSessionChannelID uint16

	writeToChannel := func(s string) {
		if sshSessionStdinChannel != nil {
			fmt.Fprint(sshSessionStdinChannel, s)
		}
	}
	simpleNewFunction("writeToChannel", false, func(this js.Value, args []js.Value) any {
		writeToChannel(args[0].String())
		return false
	})

	ctx := context.Background()

	if mp, err = NewMeepo(trackerdAddr, trackerdURL); err != nil {
		panic(err)
	}

	setElementInnerHTML("id", mp.Addr().String())
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			target := getInputValue("target")
			if target == "" {
				setElementInnerHTML("transportState", "unknown")
				setElementInnerHTML("channelState", "unknown")
				continue
			}

			targetAddr, err := addr.FromString(target)
			if err != nil {
				setElementInnerHTML("transportState", "unknown")
				setElementInnerHTML("channelState", "unknown")
				continue
			}

			t, err := mp.GetTransport(ctx, targetAddr)
			if err != nil {
				setElementInnerHTML("transportState", "unknown")
				setElementInnerHTML("channelState", "unknown")
				continue
			}

			setElementInnerHTML("transportState", t.State().String())

			if sshSessionChannelID == 0 {
				setElementInnerHTML("channelState", "unknown")
				continue
			}

			ch, err := t.GetChannel(ctx, sshSessionChannelID)
			if err != nil {
				setElementInnerHTML("channelState", "unknown")
				continue
			}

			setElementInnerHTML("channelState", ch.State().String())
		}
	}()

	doitConnect := func(target string, host string, port string, username string, password string, channelMode string) {
		defer func() {
			doitButton.Set("disabled", false)

			recovered := recover()
			if recovered != nil {
				err := recovered.(error)
				fmt.Println(err)
				return
			}

			setElementInnerHTML("doit", "disconnect")
			doitButton.Call("removeEventListener", "click", document.Underlying().Get("doitConnect"))
			doitButton.Call("removeEventListener", "click", document.Underlying().Get("doitDisconnect"))
			doitButton.Call("addEventListener", "click", document.Underlying().Get("doitDisconnect"))
		}()

		doitButton.Set("disabled", true)

		targetAddr, err := addr.FromString(target)
		if err != nil {
			panic(err)
		}

		t, err := mp.GetTransport(ctx, targetAddr)
		if err != nil {
			if !errors.Is(err, meepo_core.ErrTransportNotFound) {
				panic(err)
			}

			t, err = mp.NewTransport(ctx, targetAddr)
			if err != nil {
				panic(err)
			}
		}

		if err = t.WaitReady(); err != nil {
			panic(err)
		}

		hostAndPort := net.JoinHostPort(host, port)
		ch, err := t.NewChannel(ctx, "tcp", hostAndPort, well_known_option.WithMode(channelMode))
		if err != nil {
			panic(err)
		}

		if err = ch.WaitReady(); err != nil {
			panic(err)
		}

		conn := &myConn{ch}
		c, chans, reqs, err := ssh.NewClientConn(conn, hostAndPort, &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			panic(err)
		}

		cli := ssh.NewClient(c, chans, reqs)

		sess, err := cli.NewSession()
		if err != nil {
			panic(err)
		}

		stdout, err := sess.StdoutPipe()
		if err != nil {
			panic(err)
		}

		stderr, err := sess.StderrPipe()
		if err != nil {
			panic(err)
		}

		stdin, err := sess.StdinPipe()
		if err != nil {
			panic(err)
		}

		if err = sess.RequestPty("xterm", 80, 40, ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.ONLRET:        1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}); err != nil {
			panic(err)
		}

		if err = sess.Shell(); err != nil {
			panic(err)
		}

		go func() {
			term := document.Underlying().Get("term")
			rd := bufio.NewReader(stdout)
			for {
				b, err := rd.ReadByte()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					panic(err)
				}
				term.Call("write", string([]byte{b}))
			}
		}()

		go func() {
			term := document.Underlying().Get("term")
			rd := bufio.NewReader(stderr)
			for {
				b, err := rd.ReadByte()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					panic(err)
				}
				term.Call("write", string([]byte{b}))
			}
		}()

		sshSession = sess
		sshSessionStdinChannel = stdin
		sshSessionTarget = target
		sshSessionChannelID = ch.ID()
	}
	simpleNewFunction("doitConnect", true, func(this js.Value, args []js.Value) any {
		target := getInputValue("target")
		sshHost := getInputValue("sshHost")
		sshPort := getInputValue("sshPort")
		sshUser := getInputValue("sshUser")
		sshPassword := getInputValue("sshPassword")

		doitConnect(target, sshHost, sshPort, sshUser, sshPassword, channelMode)

		return false
	})

	doitDisconnect := func(target string, channelID uint16) {
		defer func() {
			sshSessionStdinChannel = nil
			sshSessionTarget = ""
			sshSessionChannelID = 0
			setElementInnerHTML("doit", "connect")

			doitButton.Call("removeEventListener", "click", document.Underlying().Get("doitConnect"))
			doitButton.Call("removeEventListener", "click", document.Underlying().Get("doitDisconnect"))
			doitButton.Call("addEventListener", "click", document.Underlying().Get("doitConnect"))
		}()

		if sshSessionTarget == "" || sshSessionChannelID == 0 {
			return
		}

		targetAddr, err := addr.FromString(sshSessionTarget)
		if err != nil {
			panic(err)
		}

		t, err := mp.GetTransport(ctx, targetAddr)
		if err != nil {
			panic(err)
		}

		ch, err := t.GetChannel(ctx, sshSessionChannelID)
		if err != nil {
			panic(err)
		}

		if err = sshSession.Close(); err != nil {
			if !errors.Is(err, io.EOF) {
				panic(err)
			}
		}

		if err = ch.Close(ctx); err != nil {
			panic(err)
		}
	}
	simpleNewFunction("doitDisconnect", true, func(this js.Value, args []js.Value) any {
		doitDisconnect(sshSessionTarget, sshSessionChannelID)
		return false
	})

	doitButton.Call("addEventListener", "click", document.Underlying().Get("doitConnect"))
	setElementInnerHTML("doit", "connect")

	<-q
}
