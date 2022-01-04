package meepo

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/transport"
)

const (
	METHOD_CLOSE_TRANSPORT Method = "closeTransport"
)

func (mp *Meepo) CloseTransport(peerID string) error {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "closeTransport",
		"peerID":  peerID,
	})

	tp, err := mp.getTransport(peerID)
	if err != nil {
		logger.WithError(err).Errorf("transport not found")
		return err
	}

	in := mp.createRequest(peerID, METHOD_CLOSE_TRANSPORT, nil)
	out, err := mp.doRequest(in)
	if err != nil {
		logger.WithError(err).Warningf("failed to do request")
	}

	if err = out.Err(); err != nil {
		logger.WithError(err).Warningf("failed to close transport by peer")
	}

	if err = tp.Close(); err != nil {
		logger.WithError(err).Errorf("failed to close transport")
		return err
	}

	logger.Infof("transport closed")

	return nil
}

func (mp *Meepo) onCloseTransport(dc transport.DataChannel, in packet.Packet) {
	hdr := in.Header()
	peerID := hdr.Source()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onCloseTransport",
		"peerID":  peerID,
		"session": hdr.Session(),
	})

	tp, err := mp.getTransport(peerID)
	if err != nil {
		logger.WithError(err).Errorf("transport not found")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	go func() {
		// HACK: yield cpu avoid too soon to close transport
		time.Sleep(0)

		if err = tp.Close(); err != nil {
			logger.WithError(err).Warningf("failed to close transport")
			return
		}
		logger.Tracef("transport closed")
	}()

	out := mp.createResponse(in, nil)
	mp.sendResponse(dc, out)

	logger.Infof("close transport")
}
