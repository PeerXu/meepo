package meepo

import (
	"crypto/ed25519"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/signaling"
	chain_signaling "github.com/PeerXu/meepo/pkg/signaling/chain"
	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
	mrandom "github.com/PeerXu/meepo/pkg/util/random"
	msync "github.com/PeerXu/meepo/pkg/util/sync"
)

type Meepo struct {
	rtc    *webrtc.API
	se     signaling.Engine
	socks5 Socks5Server

	transports     map[string]transport.Transport
	transportsMtx  msync.Locker
	transportLocks sync.Map

	teleportationSources map[string]*teleportation.TeleportationSource
	teleportationSinks   map[string]*teleportation.TeleportationSink
	teleportationsMtx    msync.Locker

	wireHandler    signaling.WireHandler
	wireHandlerMtx msync.Locker

	opt    ofn.Option
	logger logrus.FieldLogger

	pubk       ed25519.PublicKey
	prik       ed25519.PrivateKey
	id         string
	iceServers []string

	channelLocker msync.ChannelLocker

	broadcastCache *lru.ARCCache

	requestHandlers    map[Method]RequestHandler
	requestHandlersMtx msync.Locker

	broadcastRequestHandlers    map[Method]BroadcastRequestHandler
	broadcastRequestHandlersMtx msync.Locker

	authentication auth.Authentication
	authorization  auth.Authorization
	acl            Acl
}

func (mp *Meepo) getRawLogger() logrus.FieldLogger {
	return mp.logger
}

func (mp *Meepo) getLogger() logrus.FieldLogger {
	return mp.logger.WithFields(logrus.Fields{
		"#instance": "meepo",
		"id":        mp.GetID(),
	})
}

func (mp *Meepo) GetID() string {
	return mp.id
}

func (mp *Meepo) getICEServers() []string {
	if mp.iceServers == nil {
		mp.iceServers = cast.ToStringSlice(mp.opt.Get("iceServers").Inter())
	}

	return mp.iceServers
}

func (mp *Meepo) getWebrtcConfiguration() webrtc.Configuration {
	return webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: mp.getICEServers(),
			},
		},
	}
}

func (mp *Meepo) createRequest(dst string, meth Method, in interface{}) (out packet.Packet) {
	out, _ = packet.NewPacket(
		packet.NewHeader(generateSession(), mp.GetID(), dst, packet.Request, string(meth)),
		packet.WithData(in),
	)
	return
}

func (mp *Meepo) createBroadcastRequest(dst string, in packet.Packet) (out packet.BroadcastPacket) {
	hdr := packet.NewBroadcastHeader(generateSession(), mp.GetID(), dst, packet.BroadcastRequest, in.Header().Method(), MAX_HOP_LIMITED)
	out, _ = packet.NewBroadcastPacket(hdr, packet.WithPacket(in))
	return
}

func (mp *Meepo) createResponse(p packet.Packet, in interface{}) (out packet.Packet) {
	out, _ = packet.NewPacket(packet.InvertHeader(p.Header()), packet.WithData(in))
	return
}

func (mp *Meepo) createResponseWithError(in packet.Packet, err error) (out packet.Packet) {
	out, _ = packet.NewPacket(packet.InvertHeader(in.Header()), packet.WithError(err))
	return
}

func (mp *Meepo) createBroadcastResponse(p packet.BroadcastPacket, in packet.Packet) (out packet.BroadcastPacket) {
	hdr := packet.InvertBroadcastHeader(p.Header())
	out, _ = packet.NewBroadcastPacket(hdr, packet.WithPacket(in))
	return
}

func (mp *Meepo) createBroadcastResponseWithError(p packet.BroadcastPacket, err error) (bout packet.BroadcastPacket) {

	out, _ := packet.NewPacket(
		packet.NewHeader(p.Packet().Header().Session(), mp.GetID(), p.Packet().Header().Source(), packet.Response, p.Packet().Header().Method()),
		packet.WithError(err),
	)
	bout = mp.createBroadcastResponse(p, out)
	return
}

func (mp *Meepo) repackBroadcastRequest(dst string, p packet.BroadcastPacket) (out packet.BroadcastPacket) {
	hop := p.Header().Hop() - 1
	if hop < 0 {
		hop = 0
	}

	hdr := packet.NewBroadcastHeader(generateSession(), mp.GetID(), dst, packet.BroadcastRequest, p.Header().Method(), hop)
	out, _ = packet.NewBroadcastPacket(hdr, packet.WithPacket(p.Packet()))
	return
}

func (mp *Meepo) repackBroadcastResponse(x packet.BroadcastPacket, y packet.BroadcastPacket) (out packet.BroadcastPacket) {
	xhdr := x.Header()
	hdr := packet.NewBroadcastHeader(xhdr.Session(), mp.GetID(), xhdr.Source(), packet.BroadcastResponse, xhdr.Method(), y.Header().Hop())
	out, _ = packet.NewBroadcastPacket(hdr, packet.WithPacket(y.Packet()))
	return
}

func (mp *Meepo) addTransport(id string, tp transport.Transport) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()
	mp.addTransportNL(id, tp)
}

func (mp *Meepo) addTransportNL(id string, tp transport.Transport) {
	mp.transports[id] = tp
}

func (mp *Meepo) removeTransport(id string) {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()
	mp.removeTransportNL(id)
}

func (mp *Meepo) removeTransportNL(id string) {
	delete(mp.transports, id)
}

func (mp *Meepo) closeTeleportationsByPeerID(id string) {
	ts, _ := mp.listTeleportationsByPeerID(id)
	for _, t := range ts {
		t.Close()
	}
}

func (mp *Meepo) init() error {
	mp.initHandlers()

	return nil
}

func NewMeepo(opts ...NewMeepoOption) (*Meepo, error) {
	var logger logrus.FieldLogger
	var rtc *webrtc.API
	var se signaling.Engine
	var pubk ed25519.PublicKey
	var prik ed25519.PrivateKey
	var acl Acl
	var ok bool
	var err error

	o := newNewMeepoOption()

	for _, opt := range opts {
		opt(o)
	}

	if logger, ok = o.Get("logger").Inter().(logrus.FieldLogger); !ok {
		logger = logrus.New()
	}

	if rtc, ok = o.Get("webrtcAPI").Inter().(*webrtc.API); !ok {
		var settingEngine webrtc.SettingEngine
		settingEngine.DetachDataChannels()
		rtc = webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))
	}

	if se, ok = o.Get("signalingEngine").Inter().(signaling.Engine); !ok {
		return nil, fmt.Errorf("require signalingEngine")
	}

	if pubk, ok = o.Get("ed25519PublicKey").Inter().(ed25519.PublicKey); !ok {
		return nil, fmt.Errorf("require ed25519PublicKey")
	}

	if prik, ok = o.Get("ed25519PrivateKey").Inter().(ed25519.PrivateKey); !ok {
		return nil, fmt.Errorf("require ed25519PrivateKey")
	}

	if acl, ok = o.Get("acl").Inter().(Acl); !ok {
		return nil, fmt.Errorf("require acl")
	}

	broadcastCache, err := lru.NewARC(512)
	if err != nil {
		return nil, err
	}

	mp := &Meepo{
		rtc:                      rtc,
		se:                       se,
		pubk:                     pubk,
		prik:                     prik,
		id:                       Ed25519PublicKeyToMeepoID(pubk),
		acl:                      acl,
		transports:               make(map[string]transport.Transport),
		teleportationSources:     make(map[string]*teleportation.TeleportationSource),
		teleportationSinks:       make(map[string]*teleportation.TeleportationSink),
		requestHandlers:          make(map[Method]RequestHandler),
		broadcastRequestHandlers: make(map[Method]BroadcastRequestHandler),
		broadcastCache:           broadcastCache,
		channelLocker:            msync.NewChannelLocker(),
		opt:                      o,
		logger:                   logger,

		transportsMtx:               msync.NewLock(),
		teleportationsMtx:           msync.NewLock(),
		wireHandlerMtx:              msync.NewLock(),
		requestHandlersMtx:          msync.NewLock(),
		broadcastRequestHandlersMtx: msync.NewLock(),
	}

	mp.authentication = mp
	mp.authorization = mp

	if o.Get("asSignaling").Bool() {
		mpse := &SignalingEngineWrapper{mp}
		chse, err := signaling.NewEngine(
			"chain",
			signaling.WithLogger(logger),
			chain_signaling.WithEngine(mpse, se),
		)
		if err != nil {
			return nil, err
		}
		mp.se = chse
		mp.broadcastCache, _ = lru.NewARC(512)
	}

	if err = mp.init(); err != nil {
		return nil, err
	}

	mp.se.OnWire(mp.onNewTransport)

	return mp, nil
}

func generateSession() int32 {
	return mrandom.Random.Int31()
}
