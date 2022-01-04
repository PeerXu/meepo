package api

import (
	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/ofn"
)

type NewServerOption = ofn.OFN

func WithMeepo(meepo *meepo.Meepo) NewServerOption {
	return func(o ofn.Option) {
		o["meepo"] = meepo
	}
}
