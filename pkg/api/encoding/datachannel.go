package encoding_api

import "github.com/PeerXu/meepo/pkg/transport"

type DataChannel struct {
	Label string `json:"label"`
	State string `json:"state"`
}

func ConvertDataChannel(x transport.DataChannel) *DataChannel {
	y := &DataChannel{
		Label: x.Label(),
		State: x.State().String(),
	}

	return y
}

func ConvertDataChannels(xs []transport.DataChannel) []*DataChannel {
	var ys []*DataChannel

	for _, x := range xs {
		ys = append(ys, ConvertDataChannel(x))
	}

	return ys
}
