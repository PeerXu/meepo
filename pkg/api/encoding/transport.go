package encoding_api

import (
	"github.com/PeerXu/meepo/pkg/transport"
)

type Transport struct {
	PeerID       string         `json:"peerID"`
	State        string         `json:"state"`
	DataChannels []*DataChannel `json:"dataChannels,omitempty"`
}

func ConvertTransport(x transport.Transport) *Transport {
	dcs, _ := x.DataChannels()

	y := &Transport{
		PeerID:       x.PeerID(),
		State:        x.TransportState().String(),
		DataChannels: ConvertDataChannels(dcs),
	}

	return y
}

func ConvertTransports(xs []transport.Transport) []*Transport {
	var ys []*Transport

	for _, x := range xs {
		ys = append(ys, ConvertTransport(x))
	}

	return ys
}
