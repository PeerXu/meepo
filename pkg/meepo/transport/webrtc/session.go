package transport_webrtc

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	randomSession = Session(0)
)

type Session int32

func (s Session) String() string {
	return fmt.Sprintf("%08x", int32(s))
}

func (t *WebrtcTransport) newSession() Session {
	return newSession(t.randSrc)
}

func newSession(randSrc rand.Source) Session {
	return Session(randSrc.Int63() & (math.MaxInt32 - 1))
}

func (t *WebrtcTransport) nextSession(sess Session) Session {
	return nextSession(sess)
}

func nextSession(sess Session) Session {
	return Session(int32(sess) + 1)
}
