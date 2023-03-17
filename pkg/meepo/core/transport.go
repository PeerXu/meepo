package meepo_core

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (mp *Meepo) removeTransportNL(addr Addr) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "removeTransportNL",
		"addr":    addr.String(),
	})
	delete(mp.transports, addr)
	logger.Tracef("remove transport")
}
