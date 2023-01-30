package transport_webrtc

import (
	"context"
	"math/rand"

	"github.com/pion/webrtc/v3"

	mcontext "github.com/PeerXu/meepo/pkg/lib/context"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) registerPeerConnection(pc *webrtc.PeerConnection) Session {
	sess := t.newSession()
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "registerPeerConnection",
		"session": sess.String(),
	})

	t.peerConnections.Store(sess, pc)

	logger.Tracef("register peer connection")

	return sess
}

func (t *WebrtcTransport) unregisterPeerConnection(sess Session) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "unregisterPeerConnection",
		"session": sess.String(),
	})

	t.peerConnections.Delete(sess)

	logger.Tracef("unregister peer connection")
}

func (t *WebrtcTransport) setPeerConnection(sess Session, pc *webrtc.PeerConnection) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "setPeerConnection",
		"session": sess.String(),
	})

	t.peerConnections.Store(sess, pc)

	logger.Tracef("set peer connection")
}

func (t *WebrtcTransport) loadPeerConnection(sess Session) (*webrtc.PeerConnection, error) {
	if sess == randomSession {
		return t.loadRandomPeerConnection()
	}

	pc, ok := t.peerConnections.Load(sess)
	if !ok {
		return nil, ErrPeerConnectionNotFoundFn(sess)
	}
	return pc, nil
}

func (t *WebrtcTransport) loadRandomPeerConnection() (*webrtc.PeerConnection, error) {
	var pc *webrtc.PeerConnection

	i := 0
	r := rand.New(t.randSrc)
	t.peerConnections.Range(func(key Session, val *webrtc.PeerConnection) bool {
		if r.Float64() < (1 / (float64(i + 1))) {
			pc = val
		}
		i++
		return true
	})

	if pc == nil {
		return nil, ErrPeerConnectionNotFoundFn(randomSession)
	}

	return pc, nil
}

func (t *WebrtcTransport) loadPeerConnectionByContext(ctx context.Context) (*webrtc.PeerConnection, error) {
	sess, found := mcontext.Value[Session](ctx, OPTION_SESSION)
	if !found {
		sess = randomSession
	}
	return t.loadPeerConnection(sess)
}

func (t *WebrtcTransport) countPeerConnections() (cnt int) {
	t.peerConnections.Range(func(sess Session, pc *webrtc.PeerConnection) bool {
		cnt++
		return true
	})
	return
}

func (t *WebrtcTransport) closeAllPeerConnections() error {
	var err error
	logger := t.GetLogger().WithField("#method", "closeAllPeerConnections")

	t.peerConnections.Range(func(sess Session, pc *webrtc.PeerConnection) bool {
		if er := pc.Close(); er != nil {
			err = er
			logger.WithField("session", sess.String()).WithError(err).Debugf("failed to close peer connection")
		}
		return true
	})

	logger.Tracef("close all peer connections")

	return err
}
