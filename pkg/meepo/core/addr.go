package meepo_core

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

type Addr = meepo_interface.Addr

func (mp *Meepo) Addr() Addr {
	return mp.addr
}
