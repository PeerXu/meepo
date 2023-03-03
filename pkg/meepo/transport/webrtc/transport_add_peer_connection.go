package transport_webrtc

import (
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) addPeerConnection(sess Session) error {
	if sess == randomSession {
		sess = t.newSession()
	}

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "addPeerConnection",
		"session": sess,
	})

	pc, err := t.newPeerConnectionFunc()
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return err
	}

	if err := t.registerPeerConnection(sess, pc); err != nil {
		defer pc.Close()
		logger.WithError(err).Debugf("failed to register peer connection")
		return err
	}

	pc.OnConnectionStateChange(t.onSourceConnectionStateChange(sess))
	pc.OnDataChannel(t.onDataChannel(sess))
	go t.sourceGather(sess, t.sourceGatherFunc)

	logger.Tracef("add peer connection")

	return nil
}
