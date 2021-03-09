package meepo

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/transport"
)

type PingRequest struct {
	*Message

	Payload string `json:"payload"`
}

type PongResponse = PingRequest

func (mp *Meepo) Ping(id string, payload string) error {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "Ping",
		"peerID":  id,
	})

	req := &PingRequest{
		Message: mp.createRequest("ping"),
		Payload: payload,
	}

	out, err := mp.doRequest(id, req)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return err
	}

	res := out.(*PongResponse)
	if res.Payload != payload {
		err = fmt.Errorf("Unmatched pong payload")
		logger.WithError(err).Errorf("failed to ping")
		return err
	}

	logger.Infof("ping")

	return nil
}

func (mp *Meepo) onPing(dc transport.DataChannel, in interface{}) {
	req := in.(*PingRequest)
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onPing",
		"peerID":  req.PeerID,
		"session": req.Session,
	})

	res := &PongResponse{
		Message: mp.invertMessage(req.Message),
		Payload: req.Payload,
	}

	mp.sendMessage(dc, res)

	logger.Infof("pong")
}

func init() {
	registerDecodeMessageHelper(MESSAGE_TYPE_REQUEST, "ping", func() interface{} { return &PingRequest{} })
	registerDecodeMessageHelper(MESSAGE_TYPE_RESPONSE, "ping", func() interface{} { return &PongResponse{} })
}
