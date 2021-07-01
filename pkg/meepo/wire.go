package meepo

import (
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/transport"
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

const (
	METHOD_WIRE Method = "wire"
)

type (
	WireRequest struct {
		Descriptor *signaling.Descriptor
	}

	WireResponse struct {
		Descriptor *signaling.Descriptor
	}
)

func (mp *Meepo) Wire(peerID string, src *signaling.Descriptor) (*signaling.Descriptor, error) {
	var res WireResponse

	in := mp.createRequest(peerID, METHOD_WIRE, &WireRequest{Descriptor: src})

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "Wire",
		"session": in.Header().Session(),
	})

	out, err := mp.doBroadcastRequest(in, &doBroadcastRequestOption{
		NewGroupFunc: mgroup.NewAnyGroupFunc,
	})
	if err != nil {
		logger.WithError(err).Debugf("failed to do broadcast request")
		return nil, err
	}

	if err = out.Err(); err != nil {
		logger.WithError(err).Debugf("failed to wire")
		return nil, err
	}

	if err = out.Data(&res); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal response data")
		return nil, err
	}

	logger.Infof("wire")

	return res.Descriptor, nil
}

func (mp *Meepo) SetWireHandler(h signaling.WireHandler) {
	mp.wireHandlerMtx.Lock()
	mp.wireHandler = h
	mp.wireHandlerMtx.Unlock()
}

func (mp *Meepo) onWire(dc transport.DataChannel, bin packet.BroadcastPacket) {
	var req WireRequest
	var err error

	hdr := bin.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onWire",
		"source":  hdr.Source(),
		"session": hdr.Session(),
	})

	if err = bin.Packet().Data(&req); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal broadcast request data")
		mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
		return
	}

	mp.wireHandlerMtx.Lock()
	handler := mp.wireHandler
	mp.wireHandlerMtx.Unlock()

	if handler == nil {
		err = ErrNotWirable
		mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
		logger.WithError(err).Debugf("not wirable")
		return
	}

	desc, err := handler(req.Descriptor)
	if err != nil {
		mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
		logger.WithError(err).Debugf("failed to handle wire")
		return
	}
	logger.Tracef("handle wire")

	out := mp.createResponse(bin.Packet(), &WireResponse{Descriptor: desc})
	bout := mp.createBroadcastResponse(bin, out)
	mp.sendBroadcastResponse(dc, bout)

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
