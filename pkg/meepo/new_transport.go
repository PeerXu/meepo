package meepo

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/transport"
	_ "github.com/PeerXu/meepo/pkg/transport/loopback"
	webrtc_transport "github.com/PeerXu/meepo/pkg/transport/webrtc"
)

func (mp *Meepo) onTransportSysDataChannelCreate(dc transport.DataChannel) {
	logger := mp.getLogger().WithField("#method", "onTransportSysDataChannelCreate")

	dc.OnOpen(func() {
		innerLogger := mp.getLogger().WithFields(logrus.Fields{
			"#method": "OnOpen",
			"label":   dc.Label(),
		})
		go mp.sysDataChannelLoop(dc)
		innerLogger.Tracef("data channel opened")
	})

	logger.Tracef("data channel created")
}

func (mp *Meepo) NewTransport(peerID string) (transport.Transport, error) {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "NewTransport",
		"peerID":  peerID,
	})

	tlck := mp.getTransportLock(peerID)
	tlck.Lock()
	defer tlck.Unlock()

	if _, err = mp.getTransport(peerID); err == nil {
		err = ErrTransportExist
		logger.WithError(err).Errorf("transport already exists")
		return nil, err
	}

	var name string
	opts := []transport.NewTransportOption{
		transport.WithID(mp.GetID()),
		transport.WithPeerID(peerID),
		transport.WithLogger(mp.getRawLogger()),
	}
	if peerID == mp.GetID() {
		name = "loopback"
	} else {
		name = "webrtc"
		opts = append(opts,
			webrtc_transport.WithWebrtcAPI(mp.rtc),
			webrtc_transport.WithICEServers(mp.getICEServers()),
			webrtc_transport.AsOfferer(),
			webrtc_transport.WithOfferHook(func(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
				req, gcm, nonce, err := mp.marshalRequestDescriptor(peerID, offer)
				if err != nil {
					logger.WithError(err).Errorf("failed to marshal offer session description")
					return nil, err
				}

				res, err := mp.se.Wire(&signaling.Descriptor{ID: peerID}, req)
				if err != nil {
					logger.WithError(err).Errorf("failed to wire")
					return nil, err
				}

				answer, err := mp.unmarshalResponseDescriptor(res, gcm, nonce)
				if err != nil {
					logger.WithError(err).Errorf("failed to unmarshal answer session description")
					return nil, err
				}

				logger.Tracef("signaling engine wire")

				return answer, nil
			}),
		)
	}

	tp, err := transport.NewTransport(name, opts...)
	if err != nil {
		logger.WithError(err).Errorf("failed to new transport")
		return nil, err
	}
	tp.OnDataChannelCreate("sys", mp.onTransportSysDataChannelCreate)
	logger.Tracef("register on data channel create handler")

	h := func(transport.HandleID) {
		mp.closeTeleportationsByPeerID(peerID)
		logger.Tracef("close teleportations")

		mp.removeTransport(peerID)
		logger.Tracef("remove transport")
	}
	tp.OnTransportState(transport.TransportStateFailed, h)
	tp.OnTransportState(transport.TransportStateClosed, h)
	logger.Tracef("register on transport state change handler")

	mp.addTransport(peerID, tp)
	logger.Tracef("add transport")

	logger.Info("new transport")

	return tp, nil
}

func (mp *Meepo) onNewTransport(req *signaling.Descriptor) (*signaling.Descriptor, error) {
	var err error

	peerID := req.ID
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onNewTransport",
		"peerID":  peerID,
	})

	tlck := mp.getTransportLock(peerID)
	tlck.Lock()
	defer tlck.Unlock()

	if _, err = mp.getTransport(peerID); err == nil {
		err = ErrTransportExist
		logger.WithError(err).Errorf("transport already exists")
		return nil, err
	}

	offer, gcm, nonce, err := mp.unmarshalRequestDescriptor(req)
	if err != nil {
		return nil, err
	}

	var res *signaling.Descriptor
	var wg sync.WaitGroup
	var er error
	wg.Add(1)
	tp, err := transport.NewTransport("webrtc",
		transport.WithID(mp.GetID()),
		transport.WithPeerID(peerID),
		transport.WithLogger(mp.getRawLogger()),
		webrtc_transport.WithWebrtcAPI(mp.rtc),
		webrtc_transport.WithICEServers(mp.getICEServers()),
		webrtc_transport.WithOffer(offer),
		webrtc_transport.AsAnswerer(),
		webrtc_transport.WithAnswerHook(func(answer *webrtc.SessionDescription, hookErr error) {
			defer wg.Done()

			if hookErr != nil {
				er = hookErr
				logger.WithError(er).Errorf("failed to wire")
				return
			}

			res, er = mp.marshalResponseDescriptor(answer, gcm, nonce)
			if er != nil {
				logger.WithError(er).Errorf("failed to marshal response descriptor")
				return
			}
		}),
	)
	if err != nil {
		logger.WithError(err).Errorf("failed to new transport")
		return nil, err
	}
	tp.OnDataChannelCreate("sys", mp.onTransportSysDataChannelCreate)
	logger.Tracef("register on data channel create handler")

	h := func(transport.HandleID) {
		mp.closeTeleportationsByPeerID(peerID)
		logger.Tracef("close teleportations")

		mp.removeTransport(peerID)
		logger.Tracef("remove transport")
	}
	tp.OnTransportState(transport.TransportStateFailed, h)
	tp.OnTransportState(transport.TransportStateClosed, h)
	logger.Tracef("register on transport state change handler")

	wg.Wait()
	if er != nil {
		// logged in answer hook
		return nil, er
	}

	mp.addTransport(peerID, tp)
	logger.Tracef("add transport")

	logger.Infof("new transport")

	return res, nil
}
