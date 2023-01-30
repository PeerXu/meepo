package transport_webrtc

import (
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) sourceGather(sess Session, gather GatherFunc) {
	var err error
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "sourceGather",
	})

	defer func() {
		if err != nil {
			t.Close(t.context())
		}
	}()

	pc, err := t.loadPeerConnection(sess)
	if err != nil {
		logger.WithError(err).Debugf("failed to load peer connection")
		return
	}

	ignore, err := pc.CreateDataChannel("_IGNORE_", nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create hacked data channel")
		return
	}
	defer ignore.Close()

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create offer")
		return
	}

	gatherCompleted := webrtc.GatheringCompletePromise(pc)

	if err = pc.SetLocalDescription(offer); err != nil {
		logger.WithError(err).Debugf("failed to set offer")
		return
	}

	select {
	case <-gatherCompleted:
	case <-time.After(t.gatherTimeout):
		err = ErrGatherTimeout
		logger.WithError(err).Debugf("gather timeout")
		return
	}

	answer, err := gather(sess, *pc.LocalDescription())
	if err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	if err = pc.SetRemoteDescription(answer); err != nil {
		logger.WithError(err).Debugf("failed to set answer")
		return
	}

	logger.Tracef("gather completed")
}

func (t *WebrtcTransport) sinkGather(sess Session, offer webrtc.SessionDescription, done GatherDoneFunc) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "sinkGather",
	})

	pc, err := t.loadPeerConnection(sess)
	if err != nil {
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to load peer peer connection")
		return
	}

	if err = pc.SetRemoteDescription(offer); err != nil {
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to set offer")
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to create answer")
		return
	}

	if err = pc.SetLocalDescription(answer); err != nil {
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to set local description")
		return
	}

	select {
	case <-gatherComplete:
	case <-time.After(t.gatherTimeout):
		err = ErrGatherTimeout
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("gather timeout")
		return
	}

	answer = *pc.LocalDescription()
	if answer.Type == webrtc.SDPType(webrtc.Unknown) {
		err = ErrInvalidAnswer
		done(sess, webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	done(sess, answer, nil)
	logger.Tracef("gather completed")
}
