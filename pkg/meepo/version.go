package meepo

import "github.com/PeerXu/meepo/pkg/util/version"

func (mp *Meepo) Version() *version.V {
	return version.Get()
}
