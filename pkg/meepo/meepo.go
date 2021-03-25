package meepo

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
	"github.com/PeerXu/meepo/pkg/signaling"
	chain_signaling "github.com/PeerXu/meepo/pkg/signaling/chain"
	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
	msync "github.com/PeerXu/meepo/pkg/util/sync"
)

var (
	random *rand.Rand
)

type Meepo struct {
	rtc *webrtc.API
	se  signaling.Engine
	ae  auth.Engine

	transports    map[string]transport.Transport
	transportsMtx sync.Mutex

	teleportationSources map[string]*teleportation.TeleportationSource
	teleportationSinks   map[string]*teleportation.TeleportationSink
	teleportationsMtx    sync.Mutex

	wireHandler    signaling.WireHandler
	wireHandlerMtx sync.Mutex

	opt    objx.Map
	logger logrus.FieldLogger

	// Options
	id         string
	iceServers []string

	channelLocker msync.ChannelLocker

	broadcastCache *lru.ARCCache

	requestHandlers    map[string]RequestHandler
	requestHandlersMtx sync.Mutex

	broadcastRequestHandlers    map[string]BroadcastRequestHandler
	broadcastRequestHandlersMtx sync.Mutex
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

func (m *Message) GetMessage() *Message {
	return m
}

func (mp *Meepo) GetID() string {
	if mp.id == "" {
		mp.id = cast.ToString(mp.opt.Get("id").Inter())
	}

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

func (mp *Meepo) invertMessage(m MessageGetter) *Message {
	return InvertMessage(m.GetMessage(), mp.GetID())
}

func (mp *Meepo) invertMessageWithError(x MessageGetter, err error) *Message {
	y := InvertMessage(x.GetMessage(), mp.GetID())
	y.Error = err.Error()
	return y
}

func (mp *Meepo) invertBroadcast(b BroadcastGetter) *Broadcast {
	return InvertBroadcast(b.GetBroadcast(), mp.GetID())
}

type createRequestOption struct {
	Type string
}

func (mp *Meepo) createRequest(meth string, opts ...*createRequestOption) *Message {
	var opt *createRequestOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &createRequestOption{
			Type: MESSAGE_TYPE_REQUEST,
		}
	}

	return &Message{
		PeerID:  mp.GetID(),
		Type:    opt.Type,
		Session: generateSession(),
		Method:  meth,
	}
}

type createBroadcastOption struct {
	DetectNextHop bool
}

func (mp *Meepo) createBroadcast(destinationID string, opts ...*createBroadcastOption) *Broadcast {
	var opt *createBroadcastOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &createBroadcastOption{
			DetectNextHop: true,
		}
	}

	return &Broadcast{
		SourceID:         mp.GetID(),
		DestinationID:    destinationID,
		BroadcastSession: generateSession(),
		Hop:              MAX_HOP_LIMITED,
		DetectNextHop:    opt.DetectNextHop,
	}
}

func (mp *Meepo) createNextHopBroadcastRequest(x BroadcastRequest) BroadcastRequest {
	y := x.Copy().(BroadcastRequest)

	m := y.GetMessage()
	m.PeerID = mp.GetID()
	m.Session = generateSession()

	b := y.GetBroadcast()
	if b.Hop > MAX_HOP_LIMITED {
		b.Hop = MAX_HOP_LIMITED
	}
	b.Hop -= 1

	return y
}

func (mp *Meepo) createBroadcastResponse(out interface{}, x BroadcastRequest) Response {
	res := out.(Response).Copy().(Response)
	inverted := mp.invertMessage(x.GetMessage())
	res.GetMessage().PeerID = inverted.PeerID
	res.GetMessage().Type = inverted.Type
	res.GetMessage().Session = inverted.Session
	res.GetMessage().Method = inverted.Method

	return res
}

func (mp *Meepo) createBroadcastResponseWithError(err error, x BroadcastRequest) *BroadcastResponse {
	return &BroadcastResponse{
		Message:   mp.invertMessageWithError(x.GetMessage(), err),
		Broadcast: mp.invertBroadcast(x.GetBroadcast()),
	}
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

func (mp *Meepo) removeTeleportationsByPeerID(id string) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.removeTeleportationsByPeerIDNL(id)
}

func (mp *Meepo) removeTeleportationsByPeerIDNL(id string) {
	ts, _ := mp.listTeleportationsByPeerIDNL(id)
	for _, t := range ts {
		switch t.Portal() {
		case teleportation.PortalSink:
			mp.removeTeleportationSinkNL(t.Name())
		case teleportation.PortalSource:
			mp.removeTeleportationSourceNL(t.Name())
		}
	}
}

func (mp *Meepo) getSignatureFromUserData(ud map[string]interface{}) auth.Context {
	return cast.ToStringMap(objx.New(ud).Get("signature").Inter())
}

func (mp *Meepo) init() error {
	mp.initHandlers()

	return nil
}

func NewMeepo(opts ...NewMeepoOption) (*Meepo, error) {
	var logger logrus.FieldLogger
	var rtc *webrtc.API
	var ae auth.Engine
	var se signaling.Engine
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

	if ae, ok = o.Get("authEngine").Inter().(auth.Engine); !ok {
		return nil, fmt.Errorf("require authEngine")
	}

	if se, ok = o.Get("signalingEngine").Inter().(signaling.Engine); !ok {
		return nil, fmt.Errorf("require signalingEngine")
	}

	mp := &Meepo{
		rtc:                      rtc,
		ae:                       ae,
		se:                       se,
		transports:               make(map[string]transport.Transport),
		teleportationSources:     make(map[string]*teleportation.TeleportationSource),
		teleportationSinks:       make(map[string]*teleportation.TeleportationSink),
		requestHandlers:          make(map[string]RequestHandler),
		broadcastRequestHandlers: make(map[string]BroadcastRequestHandler),
		channelLocker:            msync.NewChannelLocker(),
		opt:                      o,
		logger:                   logger,
	}

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
	return random.Int31()
}

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
