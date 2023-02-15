package meepo_core

import (
	"context"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	tracker_interface "github.com/PeerXu/meepo/pkg/meepo/tracker/interface"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

type AddPeerConnectionHandler interface {
	OnAddPeerConnection(transport_webrtc.Session, webrtc.SessionDescription) (webrtc.SessionDescription, error)
}

func (mp *Meepo) newOnAddPeerConnectionRequest() any { return &crypto_core.Packet{} }

func (mp *Meepo) hdrOnAddPeerConnection(ctx context.Context, req any) (any, error) {
	return mp.onAddPeerConnection(req.(*crypto_core.Packet))
}

func (mp *Meepo) onAddPeerConnection(in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "onAddPeerConnection",
	})

	srcAddr, err := addr.FromBytesWithoutMagicCode(in.Source)
	if err != nil {
		logger.WithError(err).Debugf("invalid source address")
		return
	}
	logger = logger.WithField("source", srcAddr.String())

	dstAddr, err := addr.FromBytesWithoutMagicCode(in.Destination)
	if err != nil {
		logger.WithError(err).Debugf("invalid destination address")
		return
	}
	logger = logger.WithField("destination", dstAddr.String())

	if err = mp.signer.Verify(in); err != nil {
		logger.WithError(err).Debugf("failed to verify packet")
		return
	}

	if !mp.Addr().Equal(dstAddr) {
		answer, err = mp.forwardAddPeerConnectionRequest(dstAddr, in)
		if err != nil {
			logger.WithError(err).Debugf("failed to forward add peer connection request to closest trackers")
			return
		}
		logger.Tracef("forward add peer connection request")
		return
	}

	var req tracker_interface.AddPeerConnectionRequest
	if err = mp.decryptMessage(in, &req); err != nil {
		logger.WithError(err).Debugf("failed to decrypt request")
		return
	}

	done := make(chan struct{})
	defer close(done)

	t, err := mp.GetTransport(mp.context(), srcAddr)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return
	}

	h, ok := t.(AddPeerConnectionHandler)
	if !ok {
		err = ErrUnsupportedMethodFn(METHOD_ADD_PEER_CONNECTION)
		logger.WithError(err).Debugf("unsupported method")
		return
	}

	if answer, err = h.OnAddPeerConnection(transport_webrtc.Session(req.Session), req.Offer); err != nil {
		logger.WithError(err).Debugf("failed to add peer connection")
		return
	}

	logger.Tracef("on add peer connection")

	return
}

func (mp *Meepo) forwardAddPeerConnectionRequest(dstAddr addr.Addr, in *crypto_core.Packet) (answer webrtc.SessionDescription, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "forwardAddPeerConnectionRequest",
		"target":  dstAddr.String(),
	})

	out, err := mp.forwardRequest(mp.context(), dstAddr, in,
		func(tk Tracker, in *crypto_core.Packet) (any, error) { return tk.AddPeerConnection(in) },
		mp.getClosestTrackers,
		logger)
	if err != nil {
		return
	}

	answer = out.(webrtc.SessionDescription)

	return
}
