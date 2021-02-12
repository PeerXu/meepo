package meepo

import (
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
)

type Meepo struct {
	rtc *webrtc.API
	se  signaling.Engine

	opt    objx.Map
	logger logrus.FieldLogger
}

func (mp *Meepo) getLogger() logrus.FieldLogger {
	return mp.logger.WithField("id", mp.ID())
}

func (mp *Meepo) ID() string {
	return cast.ToString(mp.opt.Get("id").Inter())
}

func (mp *Meepo) iceServers() []string {
	return cast.ToStringSlice(mp.opt.Get("iceServers").Inter())
}

func (mp *Meepo) iceGatherOptions() webrtc.ICEGatherOptions {
	return webrtc.ICEGatherOptions{
		ICEServers: []webrtc.ICEServer{
			{URLs: mp.iceServers()},
		},
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
		rtc:    rtc,
		se:     se,
		opt:    o,
		logger: logger,
	}

	se.OnWire(mp.onTeleport)

	return mp, nil
}
