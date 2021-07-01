package meepo

import (
	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/transport"
	"github.com/sirupsen/logrus"
)

const (
	METHOD_CLOSE_TELEPORTATION Method = "closeTeleportation"
)

type (
	CloseTeleportationRequest struct {
		Name string
	}

	CloseTeleportationResponse struct{}
)

func (mp *Meepo) CloseTeleportation(name string) error {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "CloseTeleportation",
		"name":    name,
	})

	tp, err := mp.GetTeleportation(name, WithSourceFirst())
	if err != nil {
		logger.WithError(err).Errorf("failed to get teleportation")
		return err
	}

	in := mp.createRequest(tp.Transport().PeerID(), METHOD_CLOSE_TELEPORTATION, &CloseTeleportationRequest{Name: name})

	out, err := mp.doRequest(in)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return err
	}

	if err = out.Err(); err != nil {
		logger.WithError(err).Errorf("failed to close teleportation by peer")
		return err
	}

	go func() {
		if err = tp.Close(); err != nil {
			logger.WithError(err).Errorf("failed to close teleportation")
			return
		}
		logger.Infof("teleportation closed")
	}()

	logger.Debugf("teleportation closing")

	return nil
}

func (mp *Meepo) onCloseTeleportation(dc transport.DataChannel, in packet.Packet) {
	var err error
	var req CloseTeleportationRequest

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onCloseTeleportation",
	})

	if err = in.Data(&req); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal request data")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	logger = logger.WithField("name", req.Name)

	ts, err := mp.GetTeleportation(req.Name, WithSinkFirst())
	if err != nil {
		logger.WithError(err).Errorf("failed to get teleportation")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	go func() {
		if err = ts.Close(); err != nil {
			logger.WithError(err).Errorf("failed to close teleportation")
			return
		}
		logger.Infof("teleportation closed")
	}()

	mp.sendResponse(dc, mp.createResponse(in, &CloseTeleportationResponse{}))
	logger.Debugf("teleportation closing")
}
