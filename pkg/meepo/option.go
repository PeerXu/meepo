package meepo

import (
	"crypto/ed25519"
	"net"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/signaling"
)

func newNewMeepoOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{
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

type NewMeepoOption = ofn.OFN

func WithSignalingEngine(se signaling.Engine) ofn.OFN {
	return func(o ofn.Option) {
		o["signalingEngine"] = se
	}
}

func WithWebrtcAPI(webrtcAPI *webrtc.API) ofn.OFN {
	return func(o ofn.Option) {
		o["webrtcAPI"] = webrtcAPI
	}
}

func WithICEServers(iceServers []string) ofn.OFN {
	return func(o ofn.Option) {
		o["iceServers"] = iceServers
	}
}

func WithED25519KeyPair(pubk ed25519.PublicKey, prik ed25519.PrivateKey) ofn.OFN {
	return func(o ofn.Option) {
		o["ed25519PublicKey"] = pubk
		o["ed25519PrivateKey"] = prik
	}
}

func WithLogger(logger logrus.FieldLogger) ofn.OFN {
	return func(o ofn.Option) {
		o["logger"] = logger
	}
}

func WithGatherTimeout(d time.Duration) ofn.OFN {
	return func(o ofn.Option) {
		o["gatherTimeout"] = d
	}
}

func WithAsSignaling(b bool) ofn.OFN {
	return func(o ofn.Option) {
		o["asSignaling"] = b
	}
}

func WithAuthorizationName(name string) ofn.OFN {
	return func(o ofn.Option) {
		o["authorizationName"] = name
	}
}

func WithAuthorizationSecret(secret string) ofn.OFN {
	return func(o ofn.Option) {
		o["authorizationSecret"] = secret
	}
}

func WithAcl(acl Acl) ofn.OFN {
	return func(o ofn.Option) {
		o["acl"] = acl
	}
}

type TeleportOption = ofn.OFN

func WithLocalAddress(local net.Addr) ofn.OFN {
	return func(o ofn.Option) {
		o["local"] = local
	}
}

func WithSecret(secret string) ofn.OFN {
	return func(o ofn.Option) {
		o["secret"] = secret
	}
}

type NewTeleportationOption = ofn.OFN

func WithName(name string) ofn.OFN {
	return func(o ofn.Option) {
		o["name"] = name
	}
}

type GetTeleportationOption = ofn.OFN

func WithSourceFirst() ofn.OFN {
	return func(o ofn.Option) {
		o["getFirst"] = "source"
	}
}

func WithSinkFirst() ofn.OFN {
	return func(o ofn.Option) {
		o["getFirst"] = "sink"
	}
}

// Socks5 Server

type NewSocks5ServerOption = ofn.OFN

func WithMeepo(mp *Meepo) ofn.OFN {
	return func(o ofn.Option) {
		o["meepo"] = mp
	}
}

func WithHost(host string) ofn.OFN {
	return func(o ofn.Option) {
		o["host"] = host
	}
}

func WithPort(port int32) ofn.OFN {
	return func(o ofn.Option) {
		o["port"] = port
	}
}

// Close Transport

type CloseTransportOption = ofn.OFN

func WithGracePeriod(s string) ofn.OFN {
	return func(o ofn.Option) {
		o["gracePeriod"] = s
	}
}
