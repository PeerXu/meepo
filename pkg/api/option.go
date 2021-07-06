package api

import (
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/ofn"
)

type NewServerOption = ofn.OFN

func WithMeepo(meepo *meepo.Meepo) NewServerOption {
	return func(o objx.Map) {
		o["meepo"] = meepo
	}
}
