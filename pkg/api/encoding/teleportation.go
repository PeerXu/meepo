package encoding_api

import (
	"net"

	"github.com/PeerXu/meepo/pkg/teleportation"
)

type Addr struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

type Teleportation struct {
	Name         string         `json:"name"`
	Source       *Addr          `json:"source"`
	Sink         *Addr          `json:"sink"`
	Portal       string         `json:"portal"`
	Transport    *Transport     `json:"transport"`
	DataChannels []*DataChannel `json:"dataChannels"`
}

func ConvertAddr(x net.Addr) *Addr {
	return &Addr{
		Network: x.Network(),
		Address: x.String(),
	}
}

func ConvertTeleportation(x teleportation.Teleportation) *Teleportation {
	return &Teleportation{
		Name:   x.Name(),
		Source: ConvertAddr(x.Source()),
		Sink:   ConvertAddr(x.Sink()),
		Portal: x.Portal().String(),
		Transport: &Transport{
			PeerID: x.Transport().PeerID(),
			State:  x.Transport().TransportState().String(),
		},
		DataChannels: ConvertDataChannels(x.DataChannels()),
	}
}

func ConvertTeleportations(xs []teleportation.Teleportation) []*Teleportation {
	var ys []*Teleportation

	for _, x := range xs {
		ys = append(ys, ConvertTeleportation(x))
	}

	return ys
}
