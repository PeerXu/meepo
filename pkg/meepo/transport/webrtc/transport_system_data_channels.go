package transport_webrtc

import (
	"context"
	"io"
	"math/rand"

	"github.com/pion/webrtc/v3"

	mcontext "github.com/PeerXu/meepo/pkg/lib/context"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) registerSystemDataChannel(sess Session, dc *webrtc.DataChannel) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "registerSystemDataChannel",
		"session": sess.String(),
		"label":   dc.Label(),
	})

	t.systemDataChannels.Store(sess, dc)

	logger.Tracef("register system data channel")
}

func (t *WebrtcTransport) registerSystemReadWriteCloser(sess Session, rwc io.ReadWriteCloser) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "registerSystemReadWriteCloser",
		"session": sess.String(),
	})

	t.systemReadWriteClosers.Store(sess, rwc)

	logger.Tracef("register system rwc")
}

func (t *WebrtcTransport) unregisterSystemDataChannel(sess Session) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "unregisterSystemDataChannel",
		"session": sess.String(),
	})

	t.systemDataChannels.Delete(sess)

	logger.Tracef("unregister system data channel")
}

func (t *WebrtcTransport) unregisterSystemReadWriteCloser(sess Session) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "unregisterSystemReadWriteCloser",
		"session": sess.String(),
	})

	t.systemReadWriteClosers.Delete(sess)

	logger.Tracef("unregister system rwc")
}

func (t *WebrtcTransport) loadSystemDataChannel(sess Session) (*webrtc.DataChannel, error) {
	dc, found := t.systemDataChannels.Load(sess)
	if !found {
		return nil, ErrDataChannelNotFoundFn(sess)
	}

	return dc, nil
}

func (t *WebrtcTransport) loadRandomSystemReadWriteCloser() (io.ReadWriteCloser, error) {
	var rwc io.ReadWriteCloser

	i := 0
	r := rand.New(t.randSrc)
	t.systemReadWriteClosers.Range(func(key Session, val io.ReadWriteCloser) bool {
		if r.Float64() < (1 / float64(i+1)) {
			rwc = val
		}
		i++
		return true
	})

	if rwc == nil {
		return nil, ErrReadWriteCloserNotFoundFn(randomSession)
	}

	return rwc, nil
}

func (t *WebrtcTransport) loadSystemReadWriteCloser(sess Session) (io.ReadWriteCloser, error) {
	if sess == randomSession {
		return t.loadRandomSystemReadWriteCloser()
	}

	rwc, found := t.systemReadWriteClosers.Load(sess)
	if !found {
		return nil, ErrReadWriteCloserNotFoundFn(sess)
	}

	return rwc, nil
}

func (t *WebrtcTransport) loadSystemReadWriteCloserByContext(ctx context.Context) (io.ReadWriteCloser, error) {
	sess, found := mcontext.Value[Session](ctx, OPTION_SESSION)
	if !found {
		sess = randomSession
	}
	return t.loadSystemReadWriteCloser(sess)
}
