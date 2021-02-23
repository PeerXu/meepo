package meepo

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/transport"
	"github.com/sirupsen/logrus"
)

type CloseTeleportationRequest struct {
	Message

	Name string `json:"name"`
}

type CloseTeleportationResponse struct {
	Message
}

func (mp *Meepo) CloseTeleportation(name string) error {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	return mp.closeTeleportationNL(name)
}

func (mp *Meepo) closeTeleportationNL(name string) error {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "closeTeleportation",
		"name":    name,
	})

	tp, err := mp.getTeleportationNL(name)
	if err != nil {
		logger.WithError(err).Errorf("failed to get teleportation")
		return err
	}

	req := &CloseTeleportationRequest{
		Message: Message{
			PeerID:  mp.GetID(),
			Type:    "request",
			Session: random.Int31(),
			Method:  "closeTeleportation",
		},
		Name: name,
	}

	out, err := mp.doRequest(tp.Transport().PeerID(), req)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return err
	}
	res := out.(*CloseTeleportationResponse)
	if res.Error != "" {
		err = fmt.Errorf(res.Error)
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

func (mp *Meepo) onCloseTeleportation(dc transport.DataChannel, in interface{}) {
	var err error

	req := in.(*CloseTeleportationRequest)

	logger := mp.getLogger().WithFields(
		logrus.Fields{
			"#method": "onCloseTeleportation",
			"name":    req.Name,
		})

	ts, ok := mp.getTeleportationSink(req.Name)
	if !ok {
		err = TeleportationNotExistError
		logger.WithError(err).Errorf("failed to get teleportation sink")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}

	go func() {
		if err = ts.Close(); err != nil {
			logger.WithError(err).Errorf("failed to close teleportation")
			return
		}
		logger.Infof("teleportation closed")
	}()

	res := CloseTeleportationResponse{
		Message: mp.invertMessage(req),
	}
	mp.sendMessage(dc, &res)

	logger.Debugf("teleportation closing")
}

func init() {
	registerDecodeMessageHelper("request", "closeTeleportation", func() interface{} { return &CloseTeleportationRequest{} })
	registerDecodeMessageHelper("response", "closeTeleportation", func() interface{} { return &CloseTeleportationResponse{} })
}
