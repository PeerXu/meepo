package meepo

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/transport"
	"github.com/sirupsen/logrus"
)

type PingRequest struct {
	Message

	Payload string `json:"payload"`
}

type PongResponse = PingRequest

func (mp *Meepo) Ping(id string, payload string) error {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "Ping",
		"peerID":  id,
	})

	req := &PingRequest{
		Message: Message{
			PeerID:  mp.GetID(),
			Type:    "request",
			Session: random.Int31(),
			Method:  "ping",
		},
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
	registerDecodeMessageHelper("request", "ping", func() interface{} { return &PingRequest{} })
	registerDecodeMessageHelper("response", "ping", func() interface{} { return &PongResponse{} })
}
