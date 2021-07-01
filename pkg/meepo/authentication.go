package meepo

import (
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/ofn"
)

func WithPacket(p packet.Packet) ofn.OFN {
	return func(o objx.Map) {
		o["packet"] = p
	}
}

func (mp *Meepo) Authenticate(subject string, opts ...auth.AuthenticateOption) (err error) {
	logger := mp.getLogger().WithField("#method", "Authenticate")

	o := objx.New(map[string]interface{}{})

	for _, opt := range opts {
		opt(o)
	}

	in, ok := o.Get("packet").Inter().(packet.Packet)
	if !ok {
		logger.Debugf("require packet")
		return ErrUnauthenticated
	}

	if in.Header().Source() != subject {
		logger.Debugf("require source equal to subject")
		return ErrUnauthenticated
	}

	if err = mp.verifyPacket(in); err != nil {
		logger.WithError(err).Debugf("failed to verify packet")
		return ErrUnauthenticated
	}

	return nil
}

func WithSubject(sub string) ofn.OFN {
	return func(o objx.Map) {
		o["subject"] = sub
	}
}

type authenticatePacketOption = ofn.OFN

func (mp *Meepo) authenticatePacket(p packet.Packet, opts ...authenticatePacketOption) error {
	o := objx.New(map[string]interface{}{})
	for _, opt := range opts {
		opt(o)
	}

	sub := cast.ToString(o.Get("subject").Inter())
	if sub == "" {
		sub = p.Header().Source()
	}

	return mp.authentication.Authenticate(sub, WithPacket(p))
}
