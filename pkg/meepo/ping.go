package meepo

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/transport"
)

const (
	METHOD_PING Method = "ping"
)

type (
	PingRequest struct {
		Payload string
	}

	PongResponse struct {
		Payload string
	}
)

func (mp *Meepo) Ping(id string, payload string) error {
	var pong PongResponse

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "Ping",
		"peerID":  id,
	})

	in := mp.createRequest(id, METHOD_PING, &PingRequest{Payload: payload})

	out, err := mp.doRequest(in)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return err
	}

	if err = out.Data(&pong); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal response data")
		return err
	}

	if pong.Payload != payload {
		err = fmt.Errorf("Unmatched pong payload")
		logger.WithError(err).Errorf("failed to ping")
		return err
	}

	logger.Infof("ping")

	return nil
}

func (mp *Meepo) onPing(dc transport.DataChannel, in packet.Packet) {
	var ping PingRequest

	hdr := in.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onPing",
		"peerID":  hdr.Source(),
		"session": hdr.Session(),
	})

	if err := in.Data(&ping); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal request data")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	out := mp.createResponse(in, &PongResponse{Payload: ping.Payload})

	mp.sendResponse(dc, out)

	logger.Infof("pong")
}
