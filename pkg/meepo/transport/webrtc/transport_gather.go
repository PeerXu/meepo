package transport_webrtc

import (
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (t *WebrtcTransport) sourceGather(gather GatherFunc) {
	var err error
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "sourceGather",
	})

	defer func() {
		if err != nil {
			t.Close(t.context())
		}
	}()

	ignore, err := t.pc.CreateDataChannel("_IGNORE_", nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create hacked data channel")
		return
	}
	defer ignore.Close()

	offer, err := t.pc.CreateOffer(nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create offer")
		return
	}

	gatherCompleted := webrtc.GatheringCompletePromise(t.pc)

	if err = t.pc.SetLocalDescription(offer); err != nil {
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

	answer, err := gather(*t.pc.LocalDescription())
	if err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	if err = t.pc.SetRemoteDescription(answer); err != nil {
		logger.WithError(err).Debugf("failed to set answer")
		return
	}

	logger.Tracef("gather completed")
}

func (t *WebrtcTransport) sinkGather(offer webrtc.SessionDescription, done GatherDoneFunc) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "sinkGather",
	})

	err := t.pc.SetRemoteDescription(offer)
	if err != nil {
		done(webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to set offer")
		return
	}

	answer, err := t.pc.CreateAnswer(nil)
	if err != nil {
		done(webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to create answer")
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(t.pc)

	if err = t.pc.SetLocalDescription(answer); err != nil {
		done(webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to set local description")
		return
	}

	select {
	case <-gatherComplete:
	case <-time.After(t.gatherTimeout):
		err = ErrGatherTimeout
		done(webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("gather timeout")
		return
	}

	answer = *t.pc.LocalDescription()
	if answer.Type == webrtc.SDPType(webrtc.Unknown) {
		err = ErrInvalidAnswer
		done(webrtc.SessionDescription{}, err)
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	done(answer, nil)
	logger.Tracef("gather completed")
}
