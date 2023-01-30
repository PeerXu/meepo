package transport_webrtc

import (
	"context"
	"fmt"
	"math/rand"

	mcontext "github.com/PeerXu/meepo/pkg/lib/context"
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
	return Session(randSrc.Int63())
}

func getSessionFromContext(ctx context.Context) Session {
	sess, _ := mcontext.Value[Session](ctx, OPTION_SESSION)
	return sess
}
