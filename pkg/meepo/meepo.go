package meepo

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
)

var (
	random *rand.Rand
)

type Meepo struct {
	rtc *webrtc.API
	se  signaling.Engine

	transports           map[string]transport.Transport
	teleportationSources map[string]*teleportation.TeleportationSource
	teleportationSinks   map[string]*teleportation.TeleportationSink

	transportsMtx     sync.Mutex
	teleportationsMtx sync.Mutex

	opt    objx.Map
	logger logrus.FieldLogger

	// Options
	id         string
	iceServers []string

	sessionChannels sync.Map
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

func (mp *Meepo) invertMessage(m MessageGetter) Message {
	return InvertMessage(m.GetMessage(), mp.GetID())
}

func (mp *Meepo) invertMessageWithError(x MessageGetter, err error) Message {
	y := InvertMessage(x.GetMessage(), mp.GetID())
	y.Error = err.Error()
	return y
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

func NewMeepo(opts ...NewMeepoOption) (*Meepo, error) {
	var ok bool
	var logger logrus.FieldLogger
	var rtc *webrtc.API
	var se signaling.Engine

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

	mp := &Meepo{
		rtc:                  rtc,
		se:                   se,
		transports:           make(map[string]transport.Transport),
		teleportationSources: make(map[string]*teleportation.TeleportationSource),
		teleportationSinks:   make(map[string]*teleportation.TeleportationSink),
		opt:                  o,
		logger:               logger,
	}

	se.OnWire(mp.onNewTransport)

	return mp, nil
}

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
