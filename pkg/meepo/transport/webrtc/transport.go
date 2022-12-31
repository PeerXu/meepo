package transport_webrtc

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/xtaci/smux"

	"github.com/PeerXu/meepo/pkg/internal/dialer"
	dialer_interface "github.com/PeerXu/meepo/pkg/internal/dialer/interface"
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type WebrtcTransport struct {
	addr meepo_interface.Addr

	pc                     *webrtc.PeerConnection
	sysMtx                 lock.Locker
	sysDataChannel         *webrtc.DataChannel
	rwc                    dialer_interface.Conn
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
	closed                 atomic.Value
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

	readyErr     error
	readyTimeout time.Duration
	ready        chan struct{}
	readyCount   int32
	readyOnce    sync.Once

	fnsMtx      lock.Locker
	fns         map[string]meepo_interface.HandleFunc
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
	polls       msync.GenericsMap[*LockableChannel]

	csMtx lock.Locker
	cs    map[uint16]meepo_interface.Channel
}

func NewWebrtcSourceTransport(opts ...meepo_interface.NewTransportOption) (meepo_interface.Transport, error) {
	o := option.ApplyWithDefault(defaultNewSourceTransportOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	gather, err := GetGatherFunc(o)
	if err != nil {
		return nil, err
	}

	gatherTimeout, err := GetGatherTimeout(o)
	if err != nil {
		return nil, err
	}

	pc, err := well_known_option.GetPeerConnection(o)
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
		pc:                     pc,
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
		role:                   "source",
		currentChannelID:       0,
		readyTimeout:           readyTimeout,
		ready:                  make(chan struct{}),
		readyCount:             1,
		fnsMtx:                 lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:                    make(map[string]meepo_interface.HandleFunc),
		marshaler:              mr,
		unmarshaler:            umr,
		polls:                  msync.NewMap[*LockableChannel](),
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

	pc.OnConnectionStateChange(t.onSourceConnectionStateChange)
	pc.OnDataChannel(t.onDataChannel)
	go t.sourceGather(gather)

	return t, nil
}

func NewWebrtcSinkTransport(opts ...meepo_interface.NewTransportOption) (meepo_interface.Transport, error) {
	o := option.ApplyWithDefault(defaultNewSinkTransportOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	offer, err := GetOffer(o)
	if err != nil {
		return nil, err
	}

	gatherDone, err := GetGatherDoneFunc(o)
	if err != nil {
		return nil, err
	}

	gatherTimeout, err := GetGatherTimeout(o)
	if err != nil {
		return nil, err
	}

	pc, err := well_known_option.GetPeerConnection(o)
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
		pc:                     pc,
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
		role:                   "sink",
		currentChannelID:       1,
		readyTimeout:           readyTimeout,
		ready:                  make(chan struct{}),
		readyCount:             1,
		fnsMtx:                 lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:                    make(map[string]meepo_interface.HandleFunc),
		marshaler:              mr,
		unmarshaler:            umr,
		polls:                  msync.NewMap[*LockableChannel](),
		csMtx:                  lock.NewLock(well_known_option.WithName("tpsMtx")),
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

	pc.OnConnectionStateChange(t.onSinkConnectionStateChange)
	pc.OnDataChannel(t.onDataChannel)
	go t.sinkGather(offer, gatherDone)

	return t, nil
}

func init() {
	transport_core.RegisterNewTransportFunc("webrtc/source", NewWebrtcSourceTransport)
	transport_core.RegisterNewTransportFunc("webrtc/sink", NewWebrtcSinkTransport)
}
