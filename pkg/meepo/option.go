package meepo

import (
	"net"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
)

func newNewMeepoOption() objx.Map {
	return objx.New(map[string]interface{}{
		"iceServers": []string{
			"stun:stun.xten.com:3478",
			"stun:stun.voipbuster.com:3478",
			"stun:stun.sipgate.net:3478",
			"stun:stun.ekiga.net:3478",
			"stun:stun.ideasip.com:3478",
			"stun:stun.schlund.de:3478",
			"stun:stun.voiparound.com:3478",
			"stun:stun.voipbuster.com:3478",
			"stun:stun.voipstunt.com:3478",
			"stun:stun.counterpath.com:3478",
			"stun:stun.1und1.de:3478",
			"stun:stun.gmx.net:3478",
			"stun:stun.callwithus.com:3478",
			"stun:stun.counterpath.net:3478",
			"stun:stun.internetcalls.com:3478",
			"stun:numb.viagenie.ca:3478",
		},
		"gatherTimeout": 31 * time.Second,
	})
}

type NewMeepoOption func(objx.Map)

func WithSignalingEngine(se signaling.Engine) NewMeepoOption {
	return func(o objx.Map) {
		o["signalingEngine"] = se
	}
}

func WithWebrtcAPI(webrtcAPI *webrtc.API) NewMeepoOption {
	return func(o objx.Map) {
		o["webrtcAPI"] = webrtcAPI
	}
}

func WithICEServers(iceServers []string) NewMeepoOption {
	return func(o objx.Map) {
		o["iceServers"] = iceServers
	}
}

func WithLogger(logger logrus.FieldLogger) NewMeepoOption {
	return func(o objx.Map) {
		o["logger"] = logger
	}
}

func WithID(id string) NewMeepoOption {
	return func(o objx.Map) {
		o["id"] = id
	}
}

func WithGatherTimeout(d time.Duration) NewMeepoOption {
	return func(o objx.Map) {
		o["gatherTimeout"] = d
	}
}

type TeleportOption func(objx.Map)

func WithLocalAddress(local net.Addr) TeleportOption {
	return func(o objx.Map) {
		o["local"] = local
	}
}
