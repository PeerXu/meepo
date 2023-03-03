package transport_webrtc

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/xtaci/smux"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	matomic "github.com/PeerXu/meepo/pkg/lib/sync/atomic"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

const (
	TRANSPORT_WEBRTC_SOURCE = "webrtc/source"
	TRANSPORT_WEBRTC_SINK   = "webrtc/sink"
)

type WebrtcTransport struct {
	addr meepo_interface.Addr

	newPeerConnectionFunc NewPeerConnectionFunc
	gatherFunc            GatherFunc

	peerConnections        msync.GenericsMap[Session, *webrtc.PeerConnection]
	systemDataChannels     msync.GenericsMap[Session, *webrtc.DataChannel]
	systemReadWriteClosers msync.GenericsMap[Session, io.ReadWriteCloser]

	sysMtx                 lock.Locker
	tempDataChannelsMtx    lock.Locker
	tempDataChannels       map[string]*tempDataChannel
	tempDataChannelTimeout time.Duration
	gatherTimeout          time.Duration
	randSrc                rand.Source
	dialer                 dialer.Dialer
	onCloseCb              transport_core.OnTransportCloseFunc
	onReadyCb              transport_core.OnTransportReadyFunc
	beforeNewChannelHook   transport_core.BeforeNewChannelHook
	logger                 logging.Logger
	closed                 matomic.GenericsValue[bool]
	connectingOnce         matomic.GenericsValue[bool]
	role                   string
	currentChannelID       uint32

	muxMtx         lock.Locker
	enableMux      bool
	muxLabel       string
	muxVer         int
	muxBuf         int
	muxStreamBuf   int
	muxKeepalive   int
	muxDataChannel *webrtc.DataChannel
	muxSess        *smux.Session
	muxNocomp      bool

	kcpMtx         lock.Locker
	enableKcp      bool
	kcpLabel       string
	kcpPreset      string
	kcpCrypt       string
	kcpKey         string
	kcpMtu         int
	kcpSndwnd      int
	kcpRcvwnd      int
	kcpDataShard   int
	kcpParityShard int
	kcpDataChannel *webrtc.DataChannel
	kcpSess        *smux.Session

	readyErrVal  matomic.GenericsValue[error]
	readyTimeout time.Duration
	ready        chan struct{}
	readyCount   int32
	readyOnce    sync.Once

	fnsMtx      lock.Locker
	fns         map[string]meepo_interface.HandleFunc
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
	polls       msync.GenericsMap[string, *LockableChannel]

	csMtx lock.Locker
	cs    map[uint16]meepo_interface.Channel

	stat struct {
		failedSourceConnections int64
		failedSinkConnections   int64
	}
}

func NewWebrtcSourceTransport(opts ...meepo_interface.NewTransportOption) (meepo_interface.Transport, error) {
	o := option.ApplyWithDefault(defaultNewSourceTransportOptions(), opts...)
	t, err := newCommonWebrtcTransport(o)
	if err != nil {
		return nil, err
	}

	t.role = "source"
	t.currentChannelID = 0

	gatherOnNewFunc, err := GetGatherOnNewFunc(o)
	if err != nil {
		return nil, err
	}

	pc, err := t.newPeerConnectionFunc()
	if err != nil {
		return nil, err
	}

	sess := t.newSession()
	t.registerPeerConnection(sess, pc)

	pc.OnConnectionStateChange(t.onSourceConnectionStateChange(sess))
	pc.OnDataChannel(t.onDataChannel(sess))
	go t.sourceGather(sess, gatherOnNewFunc)

	return t, nil
}

func NewWebrtcSinkTransport(opts ...meepo_interface.NewTransportOption) (meepo_interface.Transport, error) {
	o := option.ApplyWithDefault(defaultNewSinkTransportOptions(), opts...)
	t, err := newCommonWebrtcTransport(o)
	if err != nil {
		return nil, err
	}

	t.role = "sink"
	t.currentChannelID = 1

	offer, err := GetOffer(o)
	if err != nil {
		return nil, err
	}

	gatherDoneOnNewFunc, err := GetGatherDoneOnNewFunc(o)
	if err != nil {
		return nil, err
	}

	pc, err := t.newPeerConnectionFunc()
	if err != nil {
		return nil, err
	}

	sessI32, err := GetSession(o)
	if err != nil {
		return nil, err
	}
	sess := Session(sessI32)
	t.registerPeerConnection(sess, pc)
	nextSess := t.nextSession(sess)

	pc.OnConnectionStateChange(t.onSinkConnectionStateChange(sess))
	pc.OnDataChannel(t.onDataChannel(sess))
	go t.sinkGather(sess, offer, gatherDoneOnNewFunc)
	go t.addPeerConnection(nextSess)

	return t, nil
}

func newCommonWebrtcTransport(o option.Option) (*WebrtcTransport, error) {
	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	gatherFunc, err := GetGatherFunc(o)
	if err != nil {
		return nil, err
	}

	gatherTimeout, err := GetGatherTimeout(o)
	if err != nil {
		return nil, err
	}

	newPeerConnectionFunc, err := GetNewPeerConnectionFunc(o)
	if err != nil {
		return nil, err
	}

	randSrc, err := well_known_option.GetRandSource(o)
	if err != nil {
		return nil, err
	}

	dialer, err := dialer.GetDialer(o)
	if err != nil {
		return nil, err
	}

	onTransportClose, err := transport_core.GetOnTransportCloseFunc(o)
	if err != nil {
		return nil, err
	}

	onTransportReady, err := transport_core.GetOnTransportReadyFunc(o)
	if err != nil {
		return nil, err
	}

	beforeNewChannelHook, err := transport_core.GetBeforeNewChannelHook(o)
	if err != nil {
		return nil, err
	}

	readyTimeout, err := transport_core.GetReadyTimeout(o)
	if err != nil {
		return nil, err
	}

	mr, err := marshaler.GetMarshaler(o)
	if err != nil {
		return nil, err
	}

	umr, err := marshaler.GetUnmarshaler(o)
	if err != nil {
		return nil, err
	}

	tdcTimeout, err := GetTempDataChannelTimeout(o)
	if err != nil {
		return nil, err
	}

	t := &WebrtcTransport{
		addr:                   addr,
		newPeerConnectionFunc:  newPeerConnectionFunc,
		gatherFunc:             gatherFunc,
		peerConnections:        msync.NewMap[Session, *webrtc.PeerConnection](),
		systemDataChannels:     msync.NewMap[Session, *webrtc.DataChannel](),
		systemReadWriteClosers: msync.NewMap[Session, io.ReadWriteCloser](),
		sysMtx:                 lock.NewLock(well_known_option.WithName("sysMtx")),
		tempDataChannelsMtx:    lock.NewLock(well_known_option.WithName("tempDataChannelsMtx")),
		tempDataChannels:       make(map[string]*tempDataChannel),
		tempDataChannelTimeout: tdcTimeout,
		gatherTimeout:          gatherTimeout,
		randSrc:                randSrc,
		dialer:                 dialer,
		onCloseCb:              onTransportClose,
		onReadyCb:              onTransportReady,
		beforeNewChannelHook:   beforeNewChannelHook,
		logger:                 logger,
		closed:                 matomic.NewValue[bool](),
		connectingOnce:         matomic.NewValue[bool](),
		readyErrVal:            matomic.NewValue[error](),
		readyTimeout:           readyTimeout,
		ready:                  make(chan struct{}),
		readyCount:             1,
		fnsMtx:                 lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:                    make(map[string]meepo_interface.HandleFunc),
		marshaler:              mr,
		unmarshaler:            umr,
		polls:                  msync.NewMap[string, *LockableChannel](),
		csMtx:                  lock.NewLock(well_known_option.WithName("csMtx")),
		cs:                     make(map[uint16]meepo_interface.Channel),
	}

	t.enableMux, _ = well_known_option.GetEnableMux(o)
	if t.enableMux {
		t.readyCount += 1
		t.muxMtx = lock.NewLock(well_known_option.WithName("muxMtx"))
		t.muxLabel, _ = GetMuxLabel(o)
		t.muxVer, _ = well_known_option.GetMuxVer(o)
		t.muxBuf, _ = well_known_option.GetMuxBuf(o)
		t.muxStreamBuf, _ = well_known_option.GetMuxStreamBuf(o)
		t.muxNocomp, _ = well_known_option.GetMuxNocomp(o)
	}

	t.enableKcp, _ = well_known_option.GetEnableKcp(o)
	if t.enableKcp {
		t.readyCount += 1
		t.kcpMtx = lock.NewLock(well_known_option.WithName("kcpMtx"))
		t.kcpLabel, _ = GetKcpLabel(o)
		t.kcpPreset, _ = well_known_option.GetKcpPreset(o)
		t.kcpCrypt, _ = well_known_option.GetKcpCrypt(o)
		t.kcpKey, _ = well_known_option.GetKcpKey(o)
		t.kcpMtu, _ = well_known_option.GetKcpMtu(o)
		t.kcpSndwnd, _ = well_known_option.GetKcpSndwnd(o)
		t.kcpRcvwnd, _ = well_known_option.GetKcpRcvwnd(o)
		t.kcpDataShard, _ = well_known_option.GetKcpDataShard(o)
		t.kcpParityShard, _ = well_known_option.GetKcpParityShard(o)

	}

	t.closed.Store(false)
	t.connectingOnce.Store(false)
	go t.tryCloseFailedTransport()

	return t, nil
}

func init() {
	transport_core.RegisterNewTransportFunc(TRANSPORT_WEBRTC_SOURCE, NewWebrtcSourceTransport)
	transport_core.RegisterNewTransportFunc(TRANSPORT_WEBRTC_SINK, NewWebrtcSinkTransport)
}
