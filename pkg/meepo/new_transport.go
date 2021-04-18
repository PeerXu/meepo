package meepo

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
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

	if _, err = mp.getTransport(peerID); err == nil {
		err = TransportExistError
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
				offerSignature, err := mp.ae.Sign(map[string]interface{}{
					"offer": offer,
				})

				if err != nil {
					logger.WithError(err).Errorf("failed to sign offer payload")
					return nil, err
				}

				src := &signaling.Descriptor{
					ID:                 mp.GetID(),
					SessionDescription: offer,
					UserData: map[string]interface{}{
						"signature": offerSignature,
					},
				}

				dst, err := mp.se.Wire(&signaling.Descriptor{ID: peerID}, src)
				if err != nil {
					logger.WithError(err).Errorf("failed to wire")
					return nil, err
				}

				answerSignature := mp.getSignatureFromUserData(dst.UserData)
				if err = mp.ae.Verify(map[string]interface{}{
					"answer": dst.SessionDescription,
				}, answerSignature); err != nil {
					logger.WithError(err).Errorf("failed to verify answer signature")
					return nil, err
				}

				logger.Tracef("signaling engine wire")

				return dst.SessionDescription, nil
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

func (mp *Meepo) onNewTransport(src *signaling.Descriptor) (*signaling.Descriptor, error) {
	var err error

	peerID := src.ID
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onNewTransport",
		"peerID":  peerID,
	})

	if _, err = mp.getTransport(peerID); err == nil {
		err = TransportExistError
		logger.WithError(err).Errorf("transport already exists")
		return nil, err
	}

	offerSignature := mp.getSignatureFromUserData(src.UserData)
	if err = mp.ae.Verify(map[string]interface{}{
		"offer": src.SessionDescription,
	}, offerSignature); err != nil {
		logger.WithError(err).Errorf("failed to verify offer signature")
		return nil, err
	}

	var dst signaling.Descriptor
	var wg sync.WaitGroup
	wg.Add(1)
	tp, err := transport.NewTransport("webrtc",
		transport.WithID(mp.GetID()),
		transport.WithPeerID(peerID),
		transport.WithLogger(mp.getRawLogger()),
		webrtc_transport.WithWebrtcAPI(mp.rtc),
		webrtc_transport.WithICEServers(mp.getICEServers()),
		webrtc_transport.WithOffer(src.SessionDescription),
		webrtc_transport.AsAnswerer(),
		webrtc_transport.WithAnswerHook(func(answer *webrtc.SessionDescription, hookErr error) {
			var answerSignature auth.Context

			defer wg.Done()

			if hookErr != nil {
				err = hookErr
				logger.WithError(err).Errorf("failed to wire")
				return
			}

			if answerSignature, err = mp.ae.Sign(map[string]interface{}{
				"answer": answer,
			}); err != nil {
				logger.WithError(err).Errorf("failed to sign answer payload")
				return
			}

			dst.ID = mp.GetID()
			dst.SessionDescription = answer
			dst.UserData = map[string]interface{}{
				"signature": answerSignature,
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
	if err != nil {
		// logged in answer hook
		return nil, err
	}

	mp.addTransport(peerID, tp)
	logger.Tracef("add transport")

	logger.Infof("new transport")

	return &dst, nil
}
