package transport_webrtc

import (
	"context"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

type AddPeerConnectionRequest struct {
	Session int32
	Offer   webrtc.SessionDescription
}

type AddPeerConnectionResponse struct {
	Session int32
	Answer  webrtc.SessionDescription
}

func (t *WebrtcTransport) addRemotePeerConnection(ctx context.Context, sess Session, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
	var res AddPeerConnectionResponse

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "addRemotePeerConnection",
		"session": sess.String(),
	})

	if err := t.Call(ctx, SYS_METHOD_ADD_PEER_CONNECTION, &AddPeerConnectionRequest{}, &res); err != nil {
		logger.WithError(err).Debugf("failed to add peer connection")
		return webrtc.SessionDescription{}, err
	}

	logger.Tracef("add remote peer connection")

	return res.Answer, nil
}

func (t *WebrtcTransport) onAddPeerConnection(ctx context.Context, _req any) (_res any, err error) {
	req := _req.(*AddPeerConnectionRequest)
	sess := Session(req.Session)
	res := &AddPeerConnectionResponse{Session: req.Session}

	res.Answer, err = t.addSinkPeerConnection(sess, req.Offer)
	if err != nil {
		return
	}

	return res, nil
}

func (t *WebrtcTransport) OnAddPeerConnection(sess Session, offer webrtc.SessionDescription) (answer webrtc.SessionDescription, err error) {
	return t.addSinkPeerConnection(sess, offer)
}

func (t *WebrtcTransport) addSinkPeerConnection(sess Session, offer webrtc.SessionDescription) (answer webrtc.SessionDescription, err error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "addSinkPeerConnection",
		"session": sess.String(),
	})

	pc, err := t.newPeerConnectionFunc()
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return
	}

	t.registerPeerConnection(sess, pc)
	pc.OnConnectionStateChange(t.onSinkConnectionStateChange(sess))
	pc.OnDataChannel(t.onDataChannel(sess))

	done := make(chan struct{})
	go t.sinkGather(sess, offer, func(_sess Session, _answer webrtc.SessionDescription, _err error) {
		defer close(done)
		answer = _answer
		err = _err
	})
	<-done
	if err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return
	}

	logger.Tracef("add sink peer connection")

	return
}
