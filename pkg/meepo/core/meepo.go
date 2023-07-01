package meepo_core

import (
	"math/rand"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/acl"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	crypto_interface "github.com/PeerXu/meepo/pkg/lib/crypto/interface"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/routing_table"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	meepo_routing_table_core "github.com/PeerXu/meepo/pkg/meepo/routing_table/core"
	meepo_routing_table_interface "github.com/PeerXu/meepo/pkg/meepo/routing_table/interface"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
)

type (
	Transport     = meepo_interface.Transport
	Channel       = meepo_interface.Channel
	Teleportation = meepo_interface.Teleportation
	Tracker       = tracker_interface.Tracker
)

type Meepo struct {
	addr meepo_interface.Addr

	teleportationsMtx lock.Locker
	teleportations    map[string]Teleportation

	transportsMtx lock.Locker
	transports    map[Addr]Transport

	trackersMtx     lock.Locker
	trackers        map[Addr]Tracker
	defaultTrackers map[Addr]Tracker

	routingTable          meepo_routing_table_interface.RoutingTable
	dhtAlpha              int
	poofMtx               lock.Locker
	poofInterval          time.Duration
	poofRequestCandidates int
	naviRequests          chan *NaviRequest

	acl acl.Acl

	enableMux    bool
	muxVer       int
	muxBuf       int
	muxStreamBuf int
	muxNocomp    bool

	enableKcp      bool
	kcpPreset      string
	kcpCrypt       string
	kcpKey         string
	kcpMtu         int
	kcpSndwnd      int
	kcpRcvwnd      int
	kcpDataShard   int
	kcpParityShard int

	randSrc             rand.Source
	webrtcAPI           *webrtc.API
	webrtcConfiguration webrtc.Configuration
	signer              crypto_interface.Signer
	cryptor             crypto_interface.Cryptor
	marshaler           marshaler_interface.Marshaler
	unmarshaler         marshaler_interface.Unmarshaler
	logger              logging.Logger
	closeOnce           sync.Once
	closed              chan struct{}
}

func NewMeepo(opts ...NewMeepoOption) (meepo_interface.Meepo, error) {
	o := option.ApplyWithDefault(defaultNewMeepoOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	signer, err := crypto_core.GetSigner(o)
	if err != nil {
		return nil, err
	}

	cryptor, err := crypto_core.GetCryptor(o)
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

	randSrc, err := well_known_option.GetRandSource(o)
	if err != nil {
		return nil, err
	}

	webrtcConfiguration, err := well_known_option.GetWebrtcConfiguration(o)
	if err != nil {
		return nil, err
	}

	rt, err := routing_table.NewRoutingTable(routing_table.WithID(Addr2ID(addr)))
	if err != nil {
		return nil, err
	}

	greenLine, err := meepo_routing_table_core.GetGreenLine(o)
	if err != nil {
		return nil, err
	}

	mrt := meepo_routing_table_core.NewRoutingTable(rt, greenLine)

	poofInterval, err := GetPoofInterval(o)
	if err != nil {
		return nil, err
	}

	poofRequestCandidates, err := GetPoofRequestCandidates(o)
	if err != nil {
		return nil, err
	}

	enablePoof, err := GetEnablePoof(o)
	if err != nil {
		return nil, err
	}

	defaultTrackers := make(map[Addr]Tracker)
	trackers := make(map[Addr]Tracker)
	defaultTrackersSlice, err := tracker_core.GetTrackers(o)
	if err != nil {
		return nil, err
	}
	for _, tk := range defaultTrackersSlice {
		defaultTrackers[tk.Addr()] = tk
		trackers[tk.Addr()] = tk
		rt.AddID(Addr2ID(tk.Addr())) // nolint:errcheck
	}

	dhtAlpha, err := GetDHTAlpha(o)
	if err != nil {
		return nil, err
	}

	acl, err := acl.GetAcl(o)
	if err != nil {
		return nil, err
	}

	se, err := newWebrtcSettingEngine(o)
	if err != nil {
		return nil, err
	}

	webrtcAPI := webrtc.NewAPI(webrtc.WithSettingEngine(se))

	mp := &Meepo{
		addr:                  addr,
		webrtcAPI:             webrtcAPI,
		webrtcConfiguration:   webrtcConfiguration,
		teleportationsMtx:     lock.NewLock(well_known_option.WithName("teleportationsMtx")),
		teleportations:        make(map[string]Teleportation),
		transportsMtx:         lock.NewLock(well_known_option.WithName("transportsMtx")),
		transports:            make(map[Addr]Transport),
		trackersMtx:           lock.NewLock(well_known_option.WithName("trackersMtx")),
		trackers:              trackers,
		defaultTrackers:       defaultTrackers,
		routingTable:          mrt,
		dhtAlpha:              dhtAlpha,
		poofMtx:               lock.NewLock(well_known_option.WithName("poofMtx")),
		poofInterval:          poofInterval,
		poofRequestCandidates: poofRequestCandidates,
		naviRequests:          make(chan *NaviRequest),
		acl:                   acl,
		signer:                signer,
		cryptor:               cryptor,
		randSrc:               randSrc,
		marshaler:             mr,
		unmarshaler:           umr,
		logger:                logger,
		closed:                make(chan struct{}),
	}

	mp.enableMux, _ = well_known_option.GetEnableMux(o)
	if mp.enableMux {
		mp.muxVer, _ = well_known_option.GetMuxVer(o)
		mp.muxBuf, _ = well_known_option.GetMuxBuf(o)
		mp.muxStreamBuf, _ = well_known_option.GetMuxStreamBuf(o)
		mp.muxNocomp, _ = well_known_option.GetMuxNocomp(o)
	}

	mp.enableKcp, _ = well_known_option.GetEnableKcp(o)
	if mp.enableKcp {
		mp.kcpPreset, _ = well_known_option.GetKcpPreset(o)
		mp.kcpCrypt, _ = well_known_option.GetKcpCrypt(o)
		mp.kcpKey, _ = well_known_option.GetKcpKey(o)
		mp.kcpMtu, _ = well_known_option.GetKcpMtu(o)
		mp.kcpSndwnd, _ = well_known_option.GetKcpSndwnd(o)
		mp.kcpRcvwnd, _ = well_known_option.GetKcpRcvwnd(o)
		mp.kcpDataShard, _ = well_known_option.GetKcpDataShard(o)
		mp.kcpParityShard, _ = well_known_option.GetKcpParityShard(o)
	}

	if enablePoof {
		go mp.poofLoop()
	}
	go mp.naviLoop()

	return mp, nil
}
