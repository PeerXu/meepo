package meepo

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/transport"
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

	if _, err = mp.GetTransport(peerID); err == nil {
		err = TransportExistError
		logger.WithError(err).Errorf("transport already exists")
		return nil, err
	}

	tp, err := transport.NewTransport("webrtc",
		transport.WithID(mp.GetID()),
		transport.WithPeerID(peerID),
		transport.WithLogger(mp.getRawLogger()),
		webrtc_transport.WithWebrtcAPI(mp.rtc),
		webrtc_transport.WithICEServers(mp.getICEServers()),
		webrtc_transport.AsOfferer(),
		webrtc_transport.WithOfferHook(func(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
			src := &signaling.Descriptor{
				ID:                 mp.GetID(),
				SessionDescription: offer,
			}

			dst, err := mp.se.Wire(&signaling.Descriptor{ID: peerID}, src)
			if err != nil {
				logger.WithError(err).Errorf("failed to wire")
				return nil, err
			}
			logger.Tracef("signaling engine wire")

			return dst.SessionDescription, nil
		}),
	)
	if err != nil {
		logger.WithError(err).Errorf("failed to new transport")
		return nil, err
	}
	tp.OnDataChannelCreate("sys", mp.onTransportSysDataChannelCreate)
	logger.Tracef("register on data channel create handler")

	tp.OnTransportState(transport.TransportStateFailed, func(transport.HandleID) {
		mp.removeTransport(peerID)
		logger.Tracef("remove transport")
		mp.removeTeleportationsByPeerID(peerID)
		logger.Tracef("remove teleportations")
	})
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

	if _, err = mp.GetTransport(peerID); err == nil {
		err = TransportExistError
		logger.WithError(err).Errorf("transport already exists")
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
			defer wg.Done()

			if hookErr != nil {
				err = hookErr
				logger.WithError(err).Errorf("failed to wire")
				return
			}

			dst.SessionDescription = answer
		}),
	)
	if err != nil {
		logger.WithError(err).Errorf("failed to new transport")
		return nil, err
	}
	tp.OnDataChannelCreate("sys", mp.onTransportSysDataChannelCreate)
	logger.Tracef("register on data channel create handler")

	tp.OnTransportState(transport.TransportStateFailed, func(transport.HandleID) {
		mp.removeTransport(peerID)
		logger.Tracef("remove transport")
		mp.removeTeleportationsByPeerID(peerID)
		logger.Tracef("remove teleportations")
	})
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
