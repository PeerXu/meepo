package meepo

import (
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/transport"
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

type WireRequest struct {
	*Message
	*Broadcast

	Descriptor *signaling.Descriptor `json:"descriptor"`
}

func (x *WireRequest) Copy() interface{} {
	var y WireRequest
	copier.Copy(&y, x)
	return &y
}

type WireResponse struct {
	*Message
	*Broadcast

	Descriptor *signaling.Descriptor `json:"descriptor"`
}

func (x *WireResponse) Copy() interface{} {
	var y WireResponse
	copier.Copy(&y, x)
	return &y
}

func (mp *Meepo) Wire(peerID string, src *signaling.Descriptor) (*signaling.Descriptor, error) {
	req := &WireRequest{
		Message: mp.createRequest("wire", &createRequestOption{
			Type: MESSAGE_TYPE_BROADCAST_REQUEST,
		}),
		Broadcast:  mp.createBroadcast(peerID),
		Descriptor: src,
	}

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":          "Wire",
		"session":          req.Message.Session,
		"broadcastSession": req.Broadcast.BroadcastSession,
	})

	out, err := mp.doBroadcastRequest(peerID, req, &doBroadcastRequestOption{
		NewGroupFunc: mgroup.NewAnyGroupFunc,
	})
	if err != nil {
		logger.WithError(err).Debugf("failed to do broadcast request")
		return nil, err
	}

	res := out.(*WireResponse)
	logger.Infof("wire")

	return res.Descriptor, nil
}

func (mp *Meepo) SetWireHandler(h signaling.WireHandler) {
	mp.wireHandlerMtx.Lock()
	mp.wireHandler = h
	mp.wireHandlerMtx.Unlock()
}

func (mp *Meepo) onWire(dc transport.DataChannel, in interface{}) {
	req := in.(*WireRequest)
	var err error

	logger := mp.getLogger().WithField("#method", "onWire")

	mp.wireHandlerMtx.Lock()
	handler := mp.wireHandler
	mp.wireHandlerMtx.Unlock()

	if handler == nil {
		err = NotWirableError
		res := &WireResponse{
			Message: mp.invertMessageWithError(req, err),
		}
		mp.sendMessage(dc, res)
		logger.WithError(err).Debugf("not wirable")
		return
	}

	desc, err := handler(req.Descriptor)
	if err != nil {
		res := &WireResponse{
			Message: mp.invertMessageWithError(req, err),
		}
		mp.sendMessage(dc, res)
		logger.WithError(err).Debugf("failed to handle wire")
		return
	}
	logger.Tracef("handle wire")

	res := &WireResponse{
		Message:    mp.invertMessage(req),
		Broadcast:  mp.invertBroadcast(req.Broadcast),
		Descriptor: desc,
	}
	mp.sendMessage(dc, res)

	logger.Tracef("done")
}

type SignalingEngineWrapper struct {
	meepo *Meepo
}

func (e *SignalingEngineWrapper) Wire(dst, src *signaling.Descriptor) (*signaling.Descriptor, error) {
	return e.meepo.Wire(dst.ID, src)
}

func (e *SignalingEngineWrapper) OnWire(h signaling.WireHandler) {
	e.meepo.SetWireHandler(h)
}

func (e *SignalingEngineWrapper) Close() error {
	e.meepo.SetWireHandler(nil)

	return nil
}

func init() {
	registerDecodeMessageHelper(MESSAGE_TYPE_BROADCAST_REQUEST, "wire", func() interface{} { return &WireRequest{} })
	registerDecodeMessageHelper(MESSAGE_TYPE_RESPONSE, "wire", func() interface{} { return &WireResponse{} })
}
