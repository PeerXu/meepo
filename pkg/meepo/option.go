package meepo

import (
	"net"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
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
		"gatherTimeout":       31 * time.Second,
		"waitResponseTimeout": 17 * time.Second,
	})
}

type OFN = func(objx.Map)

type NewMeepoOption = OFN

func WithSignalingEngine(se signaling.Engine) OFN {
	return func(o objx.Map) {
		o["signalingEngine"] = se
	}
}

func WithAuthEngine(ae auth.Engine) OFN {
	return func(o objx.Map) {
		o["authEngine"] = ae
	}
}

func WithWebrtcAPI(webrtcAPI *webrtc.API) OFN {
	return func(o objx.Map) {
		o["webrtcAPI"] = webrtcAPI
	}
}

func WithICEServers(iceServers []string) OFN {
	return func(o objx.Map) {
		o["iceServers"] = iceServers
	}
}

func WithLogger(logger logrus.FieldLogger) OFN {
	return func(o objx.Map) {
		o["logger"] = logger
	}
}

func WithID(id string) OFN {
	return func(o objx.Map) {
		o["id"] = id
	}
}

func WithGatherTimeout(d time.Duration) OFN {
	return func(o objx.Map) {
		o["gatherTimeout"] = d
	}
}

func WithAsSignaling(b bool) OFN {
	return func(o objx.Map) {
		o["asSignaling"] = b
	}
}

type TeleportOption = OFN

func WithLocalAddress(local net.Addr) OFN {
	return func(o objx.Map) {
		o["local"] = local
	}
}

type NewTeleportationOption = OFN

func WithName(name string) OFN {
	return func(o objx.Map) {
		o["name"] = name
	}
}
